import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:gap/gap.dart';

import '../../core/services/api_service.dart';
import '../../data/models/api_models.dart';

class OutletsPage extends ConsumerWidget {
  const OutletsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final includeDeleted = ref.watch(includeDeletedProvider);
    final outletsAsync = ref.watch(outletsProvider(includeDeleted));

    return Scaffold(
      appBar: AppBar(
        title: const Text('Workshop Outlets'),
        subtitle: Text('PostgreSQL + UUID + Soft Delete Demo'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
        actions: [
          // Toggle for including deleted items
          Switch(
            value: includeDeleted,
            onChanged: (value) {
              ref.read(includeDeletedProvider.notifier).state = value;
            },
          ),
          const Gap(8),
          Text(
            includeDeleted ? 'Include Deleted' : 'Active Only',
            style: Theme.of(context).textTheme.bodySmall,
          ),
          const Gap(16),
          // Health check button
          IconButton(
            icon: const Icon(Icons.health_and_safety),
            onPressed: () => _showHealthCheck(context, ref),
            tooltip: 'Health Check',
          ),
        ],
      ),
      body: outletsAsync.when(
        loading: () => const Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              CircularProgressIndicator(),
              Gap(16),
              Text('Loading outlets...'),
            ],
          ),
        ),
        error: (error, stackTrace) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.error_outline,
                size: 64,
                color: Theme.of(context).colorScheme.error,
              ),
              const Gap(16),
              Text(
                'Failed to load outlets',
                style: Theme.of(context).textTheme.titleLarge,
              ),
              const Gap(8),
              Text(
                error.toString(),
                style: Theme.of(context).textTheme.bodyMedium,
                textAlign: TextAlign.center,
              ),
              const Gap(16),
              ElevatedButton(
                onPressed: () => ref.refresh(outletsProvider(includeDeleted)),
                child: const Text('Retry'),
              ),
            ],
          ),
        ),
        data: (outlets) => outlets.isEmpty
            ? Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      Icons.store_outlined,
                      size: 64,
                      color: Theme.of(context).colorScheme.outline,
                    ),
                    const Gap(16),
                    Text(
                      includeDeleted
                          ? 'No outlets found'
                          : 'No active outlets found',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    const Gap(8),
                    Text(
                      includeDeleted
                          ? 'Try creating a new outlet'
                          : 'Toggle "Include Deleted" or create a new outlet',
                      style: Theme.of(context).textTheme.bodyMedium,
                    ),
                  ],
                ),
              )
            : RefreshIndicator(
                onRefresh: () async {
                  ref.refresh(outletsProvider(includeDeleted));
                },
                child: ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: outlets.length,
                  itemBuilder: (context, index) {
                    final outlet = outlets[index];
                    return OutletCard(
                      outlet: outlet,
                      onDelete: outlet.isDeleted
                          ? null
                          : () => _deleteOutlet(context, ref, outlet),
                      onRestore: outlet.isDeleted
                          ? () => _restoreOutlet(context, ref, outlet)
                          : null,
                    );
                  },
                ),
              ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _showCreateOutletDialog(context, ref),
        icon: const Icon(Icons.add),
        label: const Text('Add Outlet'),
      ),
    );
  }

  void _showHealthCheck(BuildContext context, WidgetRef ref) async {
    try {
      final apiService = ref.read(apiServiceProvider);
      final health = await apiService.healthCheck();
      
      if (context.mounted) {
        showDialog(
          context: context,
          builder: (context) => AlertDialog(
            title: const Text('Health Check'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Status: ${health['status']}'),
                const Gap(8),
                Text('Message: ${health['message']}'),
                const Gap(8),
                Text('Database: ${health['database']}'),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text('OK'),
              ),
            ],
          ),
        );
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Health check failed: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  void _deleteOutlet(BuildContext context, WidgetRef ref, OutletModel outlet) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Outlet'),
        content: Text('Are you sure you want to delete "${outlet.name}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: const Text('Delete'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      try {
        final apiService = ref.read(apiServiceProvider);
        await apiService.softDeleteOutlet(outlet.id);
        
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Outlet "${outlet.name}" deleted'),
              action: SnackBarAction(
                label: 'Undo',
                onPressed: () => _restoreOutlet(context, ref, outlet),
              ),
            ),
          );
        }
        
        // Refresh the list
        ref.refresh(outletsProvider(ref.read(includeDeletedProvider)));
      } catch (e) {
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Failed to delete outlet: $e'),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    }
  }

  void _restoreOutlet(BuildContext context, WidgetRef ref, OutletModel outlet) async {
    try {
      final apiService = ref.read(apiServiceProvider);
      await apiService.restoreOutlet(outlet.id);
      
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Outlet "${outlet.name}" restored'),
            backgroundColor: Colors.green,
          ),
        );
      }
      
      // Refresh the list
      ref.refresh(outletsProvider(ref.read(includeDeletedProvider)));
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to restore outlet: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  void _showCreateOutletDialog(BuildContext context, WidgetRef ref) {
    final nameController = TextEditingController();
    final addressController = TextEditingController();
    final phoneController = TextEditingController();
    final emailController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Create New Outlet'),
        content: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: nameController,
                decoration: const InputDecoration(
                  labelText: 'Name',
                  hintText: 'Enter outlet name',
                ),
              ),
              const Gap(16),
              TextField(
                controller: addressController,
                decoration: const InputDecoration(
                  labelText: 'Address',
                  hintText: 'Enter outlet address',
                ),
                maxLines: 2,
              ),
              const Gap(16),
              TextField(
                controller: phoneController,
                decoration: const InputDecoration(
                  labelText: 'Phone',
                  hintText: 'Enter phone number',
                ),
              ),
              const Gap(16),
              TextField(
                controller: emailController,
                decoration: const InputDecoration(
                  labelText: 'Email',
                  hintText: 'Enter email address',
                ),
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () async {
              if (nameController.text.trim().isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                    content: Text('Please enter outlet name'),
                    backgroundColor: Colors.red,
                  ),
                );
                return;
              }

              try {
                final apiService = ref.read(apiServiceProvider);
                await apiService.createOutlet(CreateOutletRequest(
                  name: nameController.text.trim(),
                  address: addressController.text.trim(),
                  phone: phoneController.text.trim(),
                  email: emailController.text.trim(),
                ));

                if (context.mounted) {
                  Navigator.of(context).pop();
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text('Outlet "${nameController.text.trim()}" created'),
                      backgroundColor: Colors.green,
                    ),
                  );
                }

                // Refresh the list
                ref.refresh(outletsProvider(ref.read(includeDeletedProvider)));
              } catch (e) {
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text('Failed to create outlet: $e'),
                      backgroundColor: Colors.red,
                    ),
                  );
                }
              }
            },
            child: const Text('Create'),
          ),
        ],
      ),
    );
  }
}

class OutletCard extends StatelessWidget {
  final OutletModel outlet;
  final VoidCallback? onDelete;
  final VoidCallback? onRestore;

  const OutletCard({
    super.key,
    required this.outlet,
    this.onDelete,
    this.onRestore,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    outlet.name,
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      decoration: outlet.isDeleted 
                          ? TextDecoration.lineThrough 
                          : null,
                      color: outlet.isDeleted 
                          ? Theme.of(context).colorScheme.onSurface.withOpacity(0.5)
                          : null,
                    ),
                  ),
                ),
                if (outlet.isDeleted)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.red.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      'DELETED',
                      style: TextStyle(
                        color: Colors.red,
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
              ],
            ),
            const Gap(8),
            Text(
              'ID: ${outlet.id}',
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                fontFamily: 'monospace',
                color: Theme.of(context).colorScheme.outline,
              ),
            ),
            const Gap(4),
            Text(
              outlet.address,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
            const Gap(4),
            Row(
              children: [
                Icon(
                  Icons.phone,
                  size: 16,
                  color: Theme.of(context).colorScheme.outline,
                ),
                const Gap(4),
                Text(
                  outlet.phone,
                  style: Theme.of(context).textTheme.bodySmall,
                ),
                const Gap(16),
                Icon(
                  Icons.email,
                  size: 16,
                  color: Theme.of(context).colorScheme.outline,
                ),
                const Gap(4),
                Expanded(
                  child: Text(
                    outlet.email,
                    style: Theme.of(context).textTheme.bodySmall,
                  ),
                ),
              ],
            ),
            const Gap(8),
            Row(
              children: [
                Text(
                  'Created: ${_formatDate(outlet.createdAt)}',
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: Theme.of(context).colorScheme.outline,
                  ),
                ),
                const Spacer(),
                if (onRestore != null)
                  TextButton.icon(
                    onPressed: onRestore,
                    icon: const Icon(Icons.restore, size: 16),
                    label: const Text('Restore'),
                    style: TextButton.styleFrom(foregroundColor: Colors.green),
                  ),
                if (onDelete != null)
                  TextButton.icon(
                    onPressed: onDelete,
                    icon: const Icon(Icons.delete, size: 16),
                    label: const Text('Delete'),
                    style: TextButton.styleFrom(foregroundColor: Colors.red),
                  ),
              ],
            ),
            if (outlet.isDeleted) ...[
              const Gap(8),
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: Colors.red.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(4),
                ),
                child: Row(
                  children: [
                    Icon(
                      Icons.delete_outline,
                      size: 16,
                      color: Colors.red,
                    ),
                    const Gap(8),
                    Expanded(
                      child: Text(
                        'Deleted: ${_formatDate(outlet.deletedAt!)}',
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: Colors.red,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year} ${date.hour}:${date.minute}';
  }
}