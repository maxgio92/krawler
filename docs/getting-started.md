# Getting started

Let's imagine you want to list the available CentOS kernel releases, scraping default mirrors. You do it by running:

```
krawler ls centos -o yaml
```

## Configuration

A configuration lets you configure parameters for the crawling, like the mirrors to scrape.

The default configuration file path is `$HOME/.krawler.yaml`. You can specify a custom path with the `--config` option.

When a configuration is not present, the [default configurations](https://github.com/maxgio92/krawler/tree/main/pkg/scrape/defaults.go) for repositories are used.

For a detailed overview see the [**reference**](/reference).

Moreover, sample configurations are available [here](https://github.com/maxgio92/krawler/tree/main/testdata).

