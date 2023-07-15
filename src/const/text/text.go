package text

const CloudSync = `Welcome to CloudSync`
const CliVersion = `CloudSync CLI version 0.0.1`

const Azure_Container_Download = `
|Examples
>download all blobs in a specific container to default location at path ~/Downloads/<myAccount>/<myContainer>
#cloudsync azure container download --account-name <myAccountName> --container <myContainer> --sasKey <mySASkey>

>download all blobs in a specific container to a specific path
#cloudsync azure container download --account-name <myAccountName> --container <myContainer> --sasKey <mySASkey> --save-to /save/to/path

>you can also specify flag input by Environment Variables starts with CLOUSCYNC_ENV_<your flag without -- and replace - by _ > like
#export CLOUSCYNC_ENV_ACCOUNT_NAME=<myAccountName>
#export CLOUSCYNC_ENV_CONTAINER=<myContainer>
#export CLOUSCYNC_ENV_SASKEY=<mySASkey>

#cloudsync azure container download --save-to /save/to/path

|Documentations
~http://github.com/epiHATR/cloudsync/examples/azure/container/download`
