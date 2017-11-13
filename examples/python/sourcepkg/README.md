This is an example of using the Fission build pipeline.


## Set up your environment

```
fission env create --name python --version 2 --image fission/python-env-2.7:0.4.0rc --builder fission/python-build-env-2.7:0.4.0rc
```

(The version in this example is 0.4.0rc, you can update it to match the release you're using.)

## Create a source archive

```
zip -jr func.zip *.py *.txt *.sh
```


## Create a function using that as the source

```
fission function create --name my-source-test --env python --src func.zip --entrypoint "user.main" --buildcmd "./build.sh" 
```

### Check on the package using kubectl

Find the package name from the function:

```
kubectl get function my-source-test -o yaml
```

(Look for `spec.package.packageRef.name` in this YAML.)

Next, find the package status:

```
kubectl get package my-source-test-XXXX -o yaml
```

You should see `buildstatus: succeeded` in the status.


### Test your function

Now you're ready to actually test your function:

```
fission fn test --name my-source-test
```
