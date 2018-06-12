# Demo 4

A few things are introduced here. Notibly, being able to use a more interesting source as an event to invoke our function (S3) as well as other important features such as environment vairables.

This demo needs a bit of setup unfortunately:

1. Run `apex init`
2. Delete the genereated `functions/hello/` directory
3. Run `apex deploy`
3. Create an S3 bucket, with an `input/` and `output/` "directory"
4. On `input/`, create an event for `ObjectCreated` that launches our newly created Lambda function
5. Note `project-example.json`. You'll need to copy over the `environment` section, and you may want to play with the memory/timout settings. Add the name of the bucket you just created
6. Inside AWS IAM, add permissions for your role to read/write to your S3 bucket
7. Redeploy with `apex deploy`

Now upload an image to `<s3bucket>/input/`. After a few seconds, check `<s3bucket>/output/` and there should be a new image in there from our function!

Clearly this is not ideal. In the real world, we could utilise `apex infra` to do most of this work, or use other tooling such as AWS SAM. But both these things are out of scope. Demo 6 has a better version of this which is much easier to understand, so if in doubt skip to that :)
