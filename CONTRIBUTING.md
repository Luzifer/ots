# Contributing to this project

Contributions are encouraged and welcome. Thank you very much in advance for taking your time to contribute!

## Code of Conduct

This project adheres to the [Contributor Covenant 3.0 Code of Conduct](https://www.contributor-covenant.org/version/3/0/).

By participating, you are expected to uphold this code. Please report unacceptable behavior to the [project maintainer](mailto:help@luzifer.io).

## How to Contribute

For small fixes, documentation improvements, translations, and other low-risk changes, feel free to open a pull request directly.

For new features, larger behavior changes or breaking changes, please open an issue first to discuss the idea before starting the implementation. This helps avoid wasted work in case the feature does not fit the project or needs a different approach.

Bug fixes are welcome as pull requests, especially when they include a clear description of the issue being fixed.

## AI-Assisted Contributions

> "A computer can never be held accountable, therefore a computer must never make a management decision."\
> – IBM Training Manual, 1979

While AI-assisted development is becoming more common, current AI systems do not truly include intelligence despite the name. The current models are merely a huge accumulation of tokens with probabilities assigned which then can be used to "predict" code.

Therefore all AI-assisted contributions must clearly state the assistance of AI tooling in an `Assisted-by` tag in every commit.

Additionally the AI tooling is not legally capable of signing off the DCO, therefore AI tooling **must not** sign-off on contributions. The required `Signed-off-by` tag must be added by the person who fully reviewed the changes, can certify the DCO applies to the change and takes full responsibility for the contribution.

Format for the tool attribution:

```
Assisted-by: AGENT_NAME:MODEL_VERSION [TOOLS USED...]
```

Complete example:

```
feat: implement feature XYZ

[optional description of the commit]

Assisted-by: ChatGPT:GPT-5 Codex
Signed-off-by: Random J Developer <random@developer.example.org>
```

Maintainers cannot reliably detect whether AI tooling was used. This requirement therefore relies on contributor honesty. If a contribution appears to be AI-assisted, maintainers may ask the contributor to confirm whether AI tooling was used and, if so, to amend the commit with the required `Assisted-by` tag.

## Developer Certificate of Origin

All commits must be signed off using `git commit -s` (i.e. `Signed-off-by: Random J Developer <random@developer.example.org>`) to indicate you've read the document and you certify your contribution meets the criteria defined in it.

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.


Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```
