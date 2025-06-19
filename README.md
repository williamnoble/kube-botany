# kube-botany

kube-botany is a web application to manage plants in your Kubernetes cluster. Its counterpart is: [kube-botany-operator](https://github.com/williamnoble/kube-botany-operator) which I'm in the process of writing. The idea is that via Kubernetes Custom Resource you can create a Plant e.g. a Sunflower, Bonsai with Motif (e.g. a Golang gopher) and overtime it grows and each day an image is generated to show its growth. Plants have watering requirements which affect the resulting image.

The web application is **partially complete**. When I started writing, I didn't realise OpenRouter does not support image
generation, and I don't have or want an OpenAI API account thus, sample images were generated in the web ui and it mocks
image generation at the moment. Most other functions are complete, with tests.

## Overview

![kube-botany overview](assets/screenshot.png)

## TODO (missing features)
- [ ] Implement Operator (currently in-progress, I'm writing with kube-builder and experimenting with controller-runtime
  directly).
- [ ] Implement Image Generation via OpenRouter, currently OpenRouter doesn't support Image Generation hence, I stopped
  developing for the time being. Because I can't experiment within the backend I haven't worked properly on the correct
  prompts, I do however, have the prompts I used in the web ui, and they seem to produce fairly consistent results.
- [ ] Improve image caching, add s3/minio fs. Also, ensure the store uses a new in-memory cache for image storing.
- [X] ~~DTO for communication between Operator/Backend needs work.~~
- [X] Consider changing `NamedspacedName` to `ID` it's clearer this works without an Operator component.


### Motivation
I had a mini-sabbatical from work, and I wanted to write some personal projects between travel. Further to this, I wanted to write a Kubernetes Operator, and this was just an idea that appealed to me :) Granted, it took me longer to write the backend, in part because I originally had overly complicated growth logic which is now removed. 