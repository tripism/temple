# temple
Directory convention based template inheritance for Golang.

  * Each directory is a template
  * Parent templates and partials are inherited and available to children
  * Special `base.temple` will be called first
  * Use your own file extensions (as long as `.temple` appears)