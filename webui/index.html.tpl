<!doctype html>
<html>
  <head>
    <title>ACH File Viewer</title>
  </head>
  <body>
    {{ range $source := .Sources }}
    <strong>{{ $source.ID }}</strong> ({{ $source.Type }})
    {{ range $file := $source.Files }}
    <br /><a href="{{ $file.Path }}">{{ $file.Filename }}</a>
    {{ end }}
    <br /></br />
    {{ end }}
  </body>
</html>
