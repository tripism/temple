# temple
Convention based template abstraction for Golang.

  * Each directory is a template made up of many partials
  * Ancestors are inherited and available to children
  * Partials are available to every template
  * Auto-reload: Recompiles templates whenever a dependent file changes
  * Name a file `base.temple` to create master templates
  * Use your own file extensions (as long as `.temple` appears)

## Get started

Create a file structure that looks something like this:

```
/templates
  /home                 // tpl.Get("home")
    /welcome            // tpl.Get("home.welcome")
      content.temple
  /about                // tpl.Get("about")
    content.temple
  base.temple           // base file for everything
  copyright.temple      // partial available to base, welcome and about
  /_punctuation
    period.temple       // shared component
```

Insert the following 

#### `base.temple`

```
<div>
{{ template "content" . }}
</div>
```

#### `home/welcome/content.temple`

```
Welcome to the site {{ template "copyright" }}
```

#### `about/content.temple`

```
{{ template "copyright" }} All about the site{{ template "punctuation.period" }}
```

#### `copyright.temple`

```
(Copyright &copy; 2015)
```

### Using temple

First, process the files by calling `temple.New`:

```
tpl, err := temple.New("/path/to/templates")
if err != nil {
	log.Fatalln("fatal error:", err)
}
```

Then use the `tpl` as a map, using dot notation to access the
templates, and call the `Execute` method:

```
http.HandleFunc("/home/welcome", func(w http.ResponseWriter, r *http.Request) {
	if err := tpl.Get("home.welcome").Execute(w, nil); err != nil {
    // TODO: handle error
  }
})
```

### Handling functions

Since Temple pre-preocesses tempaltes, you need to specify any functions up-front via the `temple.NewFuncs` method call.

## Using the reloader

Use `temple.New` as normal:

```
tpl, err := temple.New("/path/to/templates")
if err != nil {
  log.Fatalln("fatal error:", err)
}
```

Then mix in a `temple/reloader`:

```
rl, err := reloader.New(tpl)
if err != nil {
  log.Fatalln("fatal error:", err)
}
defer rl.Close() // be sure to close it
```

Now whenever a template file changes, the reloader will call the `Reload` method.