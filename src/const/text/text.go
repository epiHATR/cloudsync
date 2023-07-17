package text

const CloudSync = `|Cloudsync CLI by Hai Tran <hidetran@gmail.com>

|The cloudsync CLI supports methods that help you download, copy, and synchronize resources such as blobs and bucket files between cloud providers.

|Providers				Object types                Methods
>Azure					blob                        download
>AWS					container                   copy 
>					bucket                      clone
>					file-share                  upload

|Documentations
~https://github.com/epiHATR/cloudsync/tree/main/docs`

const Azure_Container_Download_HelpText = `
|Examples
>download all blobs in a specific container to default location at path ~/Downloads/<myAccount>/<myContainer>
#cloudsync azure container download --account-name <myAccountName> --container <myContainer> --key <myStorageAccountKey>

>download all blobs in a specific container to a specific path
#cloudsync azure container download --account-name <myAccountName> --container <myContainer> --key <myStorageAccountKey> --save-to /save/to/path

>you can also specify flag input by Environment Variables starts with CLOUSCYNC_ENV_<your flag without -- and replace - by _ > like
#export CLOUSCYNC_ENV_ACCOUNT_NAME=<myAccountName>
#export CLOUSCYNC_ENV_CONTAINER=<myContainer>
#export CLOUSCYNC_ENV_KEY=<myStorageAccountKey>

#cloudsync azure container download --save-to /save/to/path

|Documentations
~https://github.com/epiHATR/cloudsync/tree/main/docs/azure/container`

const Azure_Container_Copy_HelpText = `
|Examples
#copy container myContainer from Azure storage account myStorageAccount1 to storage account myStorageAccount2 with storage account key
>cloudsync azure container copy --source-account myStorageAccount1 --source-container myContainer --source-key myAccount1Key --destination-account myStorageAccount2 --destination-key myAccount2Key

#copy container myContainer from Azure storage account myStorageAccount1 to storage account myStorageAccount2 with storage account connection string
>cloudsync azure container copy -source-container myContainer --source-connection-string myAccount1Conn --destination-account myStorageAccount2 --destination-connection-string myAccount2Conn

|Documentations
~https://github.com/epiHATR/cloudsync/tree/main/docs/azure/container`
