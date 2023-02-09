# Create a custom package

To create a custom package, one can begin by running the `package generate` command of the CLI. The resulting output will have folder structure with files as below:

```
.
├── docker-compose.dev.yml
├── docker-compose.yml
├── package-metadata.json
└── swarm.sh
```

{% hint style="info" %}
Review packages in [https://github.com/jembi/platform](https://github.com/jembi/platform), for examples on how to structure packages for importing configs, using several convenience utilities, running in clustered mode, and more.
{% endhint %}

