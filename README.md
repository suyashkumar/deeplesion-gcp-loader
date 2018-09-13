# DeepLesion GCP Loader
This program is a simple way to fetch, uncompress, and upload the [DeepLesion](https://www.nih.gov/news-events/news-releases/nih-clinical-center-releases-dataset-32000-ct-images) dataset of 32,000 CT images into a google cloud bucket. Usage is simple:

```sh
./deeplesion-loader --removeFiles=true --bucketName=my-bucket
```
Will download each 4GB zip from the dataset, unzip it, and upload the images to `my-bucket`. This configuration with `removeFiles=true` will delete each zip file after it has successfully uploaded the contents to GCP. 

```sh
./deeplesion-loader --parallel=true
```
Will run all file downloads and uploads in parallel--this is *much faster*, but requires more disk space and resources. 

**Note:** You must ensure the machine running this program has write access to your GCP bucket (or that GCP application deafult credentials are set). See the section below for more details

## General installation and setup
You can simply download the right binary from the [releases tab](https://github.com/suyashkumar/deeplesion-gcp-loader/releases) and run it like detailed above. You can also fetch the binary from the commandline using the following command:

```sh
wget -qO- $BINARY_RELEASE_LINK | tar xvz
```

where `$BINARY_RELEASE_LINK` is the link of the download from the [releases tab](https://github.com/suyashkumar/deeplesion-gcp-loader/releases).

### Ensuring GCP Write access
The machine this program runs on needs to have write access to your bucket. This can be done in two ways:
* Ensure [application default credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login) are set. Usually: `gcloud auth application-default login` will do it 
* Or you can spin up a GCP virtual machine that has the "Storage" API permission set to "Read Write" which can be done when creating the VM by clicking "Set access for each API"
