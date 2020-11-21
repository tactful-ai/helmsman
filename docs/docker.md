Robban
===

Visual Web UI to monitor multiple helm releases and to promote images accross them.

# How to use?


```console
docker run -v /path/to/config/dir:/opt/config \
    -p 8080:8080 tactful/robban \
    -f config/myconfig.yaml
```

Then visit http://localhost:8080 to manage your helm releases

## Configuration

Refer to [Documentation](https://github.com/tactful-ai/robban/#documentation) for guides on how to configure your cluster.

# License

This Image and the original repo are under [MIT License](https://github.com/tactful-ai/robban/blob/master/LICENSE).

# User Feedback

## Issues

If you have any problems with or questions about this image, please contact us through a [GitHub issue](https://github.com/SISheogorath/readme-to-dockerhub/issues).


## Contributing

You are invited to contribute new features, fixes, or updates, large or small; we are always thrilled to receive pull requests.
Refer to Robban [contribution guide](https://github.com/tactful-ai/robban/blob/master/CONTRIBUTION.md) for more info.  
