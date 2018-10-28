#!/bin/bash

curl -X DELETE https://pixe.la/v1/users/jagijagijag1/graphs/garmin-stress/20180719 -H 'X-USER-TOKEN:Mk6LGS5b'
aws s3 cp "./2018-10-26 18.29.25.png" s3://pixela-datasource-stress-img-bucket/