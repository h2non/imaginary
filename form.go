package main

const formText = `
<html>
<body>
<h1>Resize</h1>
<form method="POST" action="/resize?width=300&height=200&type=png" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Crop</h1>
<form method="POST" action="/crop?width=562&height=562&quality=95" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Extract</h1>
<form method="POST" action="/extract?top=100&left=100&areawidth=300&areaheight=150&quality=95" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Enlarge</h1>
<form method="POST" action="/enlarge?width=1440&height=900&quality=95" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Flip</h1>
<form method="POST" action="/flip" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Flop</h1>
<form method="POST" action="/flop" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Rotate (180)</h1>
<form method="POST" action="/rotate?rotate=180" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Thumbnail</h1>
<form method="POST" action="/thumbnail?width=100" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Zoom</h1>
<form method="POST" action="/zoom?factor=2&areawidth=300&top=80&left=80" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Watermark</h1>
<form method="POST" action="/watermark?width=100&text=Hello&font=sans%2012&opacity=0.5&color=255,200,50" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Convert to PNG</h1>
<form method="POST" action="/convert?type=png" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
<h1>Info (image metadata)</h1>
<form method="POST" action="/info" enctype="multipart/form-data">
  <input type="file" name="file" />
  <input type="submit" value="Upload" />
</form>
</body>
</html>
`
