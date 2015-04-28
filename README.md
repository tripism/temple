# temple
Directory convention based template inheritance for Golang.

  * Each directory is a template made up of many partials
  * Parent partials are inherited and available to children
  * Name a file `base.temple` to create master templates
  * Use your own file extensions (as long as `.temple` appears)

## Get started

Create a file structure that looks something like this:

```
/templates
  /home                 // tpl["home"]
    /welcome            // tpl["home.welcome"]
      content.temple
  /about                // tpl["about"]
    content.temple
  base.temple           // base file for everything
  copyright.temple      // partial available to base, welcome and about
```

Insert the following 

#### `base.temple`

```
<div>
{{ template "content" . }}
</div>
```

#### `welcome/content.temple`

```
Welcome to the site {{ template "copyright" }}
```

#### `about/content.temple`

```
{{ template "copyright" }} All about the site
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
	if err := tpl["home.welcome"].Execute(w, nil); err != nil {
    // TODO: handle error
  }
})
```
