# aws_fileserver
Basic File Server Written In Go

This code is a simple server which serves as an example of how to use go as a file server.  Static files are uploaded to AWS S3. Go will sync the local server with what's on S3 (only if the file has been changed, a new file was uploaded, or a file was removed).  Those files can then be downloaded via the server. 

Code assumes you have the credentials file setup in you .aws folder.

Possible modifications for those interested. Download directly from S3 and skip the file syncing step.. All that you would need to do is modify the src/downloader/downloader.go file. Pass the S3 key file name to the server's front-end.  You have the basic downloader code for all of that.  Use the key file name and download the S3 file either to server (and remove after download), memory (beware large files), or chunk it in real time, and you've skipped the file sync step and replaced it with a ligher version which just syncs key names.
