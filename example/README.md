# Try it out:

This folder contains an example implementation of csi-driver-sshfs.
All you need to do is follow the steps below.

> Watch out (!): This Project is for testing pursposes only, under no circumstances should it be run
    on production. This example project is for local testing only (!) and is **not monitored for security issues**.

## Prequesites
To get started, make sure you have the follwing installed:

- [MiniKube](https://minikube.sigs.k8s.io/docs/start/) or other
- [Docker](https://docs.docker.com/get-docker/)

## Build

1. Create Namespace:
    Let's first create an example Namespace called `sshfs-example`.

    ```
    $ kubectl create namespace sshfs-example
    ```

2. Verify the new namespace
    ```
    $ kubectl get namespaces
    ```

    Which should output something like:
    ```
    NAME                   STATUS   AGE
    default                Active   17d
    sshfs-example          Active   21s
    ...
    ```

3. Generate RSA Keys:
    ```
    $ ssh-keygen -f ./id_rsa -t rsa -b 4096
    $ kubectl create secret generic my-ssh-key --from-file=id_rsa=id_rsa --namespace=sshfs-example
    ```

> *Tipp:* To be extra sure you might want to delete the `id_rsa` file from your local device

4. Finally, deploy it
    ```
    $ kubectl create -f deployment.yaml --namespace=sshfs-example
    ```




