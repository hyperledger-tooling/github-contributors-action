# GitHub Contributors Action

The action can be used to fetch all contributors' information for either 
a particular repository or an organization. In case of an organization
all the repositories visible to the token user will be pulled. Resulting
contributor information is checked for duplications.

To use the action, define the template file, let the
tool get you the result in the expected format. The output from
the tool is captured in an output file. Highly useful if you would
like to highlight contributors of a project.

## Usage

The tool works by accepting following environment variables

| ENV Variable                     | What is it used for?                                                               | Default Value                                  |
|----------------------------------|------------------------------------------------------------------------------------|------------------------------------------------|
| GITHUB_AUTH_TOKEN                | GitHub access token with access to the repository                                  |                                                |
| SOURCE_GITHUB_REPOSITORY         | Repository to poll the collaborators information. It can be a simple organization field, in which case all repositories are polled for the required information                                 | hyperledger-tooling/github-contributors-action |
| CONTRIBUTORS_SECTION_PATTERN     | Pattern after which the contributors information is to be dumped                   | # # Contributors                                |
| CONTRIBUTORS_SECTION_END_PATTERN | Pattern to end adding the contributors information. This would mean all contributors information will be added in between `CONTRIBUTORS_SECTION_PATTERN` and `CONTRIBUTORS_SECTION_END_PATTERN` | # # Contributions                              |
| INPUT_TEMPLATE_FILE              | File used for creating the pattern in the output file                              | `assets/minimal.md`                            |
| FILE_WITH_PATTERN                | File where the generated data is dumped. Also, the file where pattern can be found | `README.md`                                    |

**Note 1:** An extra space is added between `##` in the `Default Value` column.
This is to avoid the action from corrupt formatting of this file.

**Note 2:** An example of pulling list of all organization contributors is
available at [Hyperledger Tooling Contributors](./all-contributors.md) page.

## Contributors
 <a href="https://github.com/arsulegai"><img src="https://avatars.githubusercontent.com/u/27664223?v=4" width="32" height="32" alt="27664223"></a>  <a href="https://github.com/hyperledger-bot"><img src="https://avatars.githubusercontent.com/u/76175814?v=4" width="32" height="32" alt="76175814"></a>  <a href="https://github.com/nidhi-singh02"><img src="https://avatars.githubusercontent.com/u/38173192?v=4" width="32" height="32" alt="38173192"></a>  <a href="https://github.com/ryjones"><img src="https://avatars.githubusercontent.com/u/466142?v=4" width="32" height="32" alt="466142"></a> 

## Contributions

Each commit must include `Signed-off-by:`
in the commit message (run `git commit -s` to auto-sign).
This sign off means you agree the commit satisfies the
[Developer Certificate of Origin(DCO)](https://developercertificate.org/).

### Build

`Makefile` is provided to ease up the development, and auto-perform lint
checks as well as run tests if any. This is to maintain a good quality on
the codebase.

Run the following command from your terminal

```shell
make
```

In case of build failure because of format check, run the following command

```shell
make format
```

### Version

Increment version by raising a PR against [VERSION](./VERSION) file.
Use the same tag while creating a release/tag.
