# GarminStress2Pixela
GarminStress2Pixela is a serverless tool to export [Garmin Connect](https://connect.garmin.com/ja-JP/) stress score to your [Pixela](https://pixe.la/) graph.

![screen](./docs/GarminStress2Pixela.png "screen")

#### Note: The tool is tailored for the image of Garmin Connect 4.12.0.14 on iPhone 7(iOS 12.01). For other environment, you shoud do tuning some part of code & parameters (e.g. position of target values in `garmin-stress2pixela/main.go`).

## Project setup
### Requirements
- Go environment
- serverless framework

### compile & deploy
```bash
git clone https://github.com/jagijagijag1/GarminStress2Pixela
cd GarminStress2Pixela
```

Describe your s3 bucket info & pixela info to `serverless.yml`
```yaml:serverless.yml
...
functions:
  garmin-stress2pixela:
    handler: bin/garmin-stress2pixela
    events:
      - s3: <bucket-name>
    # you need to fill the followings with your own
    environment:
      PIXELA_USER: <user-id>
      PIXELA_TOKEN: <your-token>
      PIXELA_GRAPH: <your-graph-id-1>
    timeout: 10
```

Then, run the following.

```bash
make
sls deploy
```
