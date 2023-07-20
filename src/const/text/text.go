package text

const CloudSync = `|Cloudsync CLI by Hai Tran <hidetran@gmail.com>

|The cloudsync CLI supports methods that help you download, copy, and synchronize resources such as blobs and bucket files between cloud providers.

|Providers				Object types                Methods
>Azure					blob                        download
>AWS					container                   copy 
>					bucket                      upload
>					file-share                  delete
>						                    list

|Documentations
~https://github.com/epiHATR/cloudsync/tree/main/docs`

const Azure_Container_Download_HelpText = `|Examples

#download all blobs in a specific container to default location at path ~/Downloads/myAccountName/myContainer
>cloudsync azure container download --account myAccountName --container myContainer --key myStorageAccountKey

#download all blobs in a specific container to a specific path
>cloudsync azure container download --account myAccountName \
>                                   --container myContainer \
>                                   --key myStorageAccountKey \
>                                   --save-to /save/to/path

#download a blob in container myContainer of storage account myAccountName to a path
>cloudsync azure container download --account myAccountName --container myContainer --blob /my/path/toBlob --save-to /my/local/path

#you can also specify flag input by Environment Variables starts with CLOUSCYNC_ENV_<your flag without -- and replace - by _ > like
>export CLOUSCYNC_ENV_ACCOUNT=myAccountName
>export CLOUSCYNC_ENV_CONTAINER=myContainer
>export CLOUSCYNC_ENV_KEY=myStorageAccountKey

>cloudsync azure container download --save-to /save/to/path

|Documentations
~https://github.com/epiHATR/cloudsync/tree/main/docs/azure/container`

const Azure_Container_Copy_HelpText = `|Examples

##copy container myContainer from Azure storage account myStorageAccount1 to storage account myStorageAccount2 with storage account key
>cloudsync azure container copy --source-account myStorageAccount1 \
>                               --source-container myContainer \
>                               --source-key myAccount1Key \
>                               --destination-account myStorageAccount2 \
>                               --destination-key myAccount2Key

##copy container myContainer from Azure storage account myStorageAccount1 to storage account myStorageAccount2 with storage account connection string
>cloudsync azure container copy --source-container myContainer \
>                               --source-connection-string myAccount1Conn \ 
>                               --destination-account myStorageAccount2 \
>                               --destination-connection-string myAccount2Conn

|Documentations
~https://github.com/epiHATR/cloudsync/tree/main/docs/azure/container`

const Azure_Container_Upload_HelpText = `|Examples

#Upload a folder at path /my/path to a storage container called myContainer of storage account myAccountName with key myAccountKey
>cloudsync azure container upload --account myAccountName --key myAccountKey --container myContainer --path /my/path

#Upload a file in a folder to a storage container called myContainer of storage account myAccountName with key myAccountKey
>cloudsync azure container upload --account myAccountName --key myAccountKey --container myContainer --path /my/path/to/file.txt

#Upload a folder at path /my/path to a storage container called myContainer with storage account connection string
>cloudsync azure container upload --container myContainer --connection-string myConnectionString --path /my/path

|Documentations
~https://github.com/epiHATR/cloudsync/tree/main/docs/azure/container`

const Base64_HelpText = `|Examples`

const Azure_Container_Delete_HelpText = `|Examples`
