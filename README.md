# AWS inspector assessment at ec2 launch

You will need the amazon inspector agent to be installed on the instance before assessment can begin.There are 2 ways to do that.

1. Using SSM agent
2. Using user data

You can find detailed steps here: https://aws.amazon.com/blogs/security/how-to-simplify-security-assessment-setup-using-ec2-systems-manager-and-amazon-inspector/


We can use "user data" to test it out for a single instance

```
#!/bin/bash

cd /tmp
curl -O https://d1wk0tztpsntt1.cloudfront.net/linux/latest/install
chmod +x /tmp/install
/tmp/install
```