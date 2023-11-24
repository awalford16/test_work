# Config Map Watcher

For eductional purposes

Trying to understand how errgroup works if subprocesses are created from informer triggers

This will create a worker go routine when a config map with a particular label is added to a cluster. If that particular configmap is deleted then it will stop the worker process.

```
Starting routine for ConfigMap: my-config-1
Starting routine for ConfigMap: my-config-2
Routine for ConfigMap my-config-2 is running...
Routine for ConfigMap my-config-1 is running...
removing capacity
Routine for ConfigMap my-config-1 stopped.
Routine for ConfigMap my-config-2 is running...
Routine for ConfigMap my-config-2 is running...
Routine for ConfigMap my-config-2 is running...
Routine for ConfigMap my-config-2 is running...
removing capacity
Routine for ConfigMap my-config-2 stopped.
```