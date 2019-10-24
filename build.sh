docker build -t milideluxe/s3-imaginary -f prod.Dockerfile .
#docker push milideluxe/s3-imaginary
rm -f imaginary.zip
rm -rf deployment
mkdir -p deployment/config
cp credentials deployment/config/
cp Dockerrun.aws.json deployment/
cd deployment
zip -r ../imaginary.zip .
echo https://eu-central-1.console.aws.amazon.com/elasticbeanstalk/home?region=eu-central-1#/environment/dashboard?applicationName=estudy-2.0&environmentId=e-eradsmm6fp