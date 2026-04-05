# aurl - Turn APIs Into Simple Commands

[![Download aurl](https://img.shields.io/badge/Download-Release%20Page-blue?style=for-the-badge)](https://github.com/ma591/aurl/releases)

## 🧭 What aurl does

aurl is a command line tool for Windows that turns an API into a command you can run from a terminal.

It works with:

- OpenAPI 3.0
- OpenAPI 3.1
- Swagger 2.0
- GraphQL

aurl can:

- detect login settings from the API
- check requests against the API spec
- build help text from the API schema
- generate docs from GraphQL introspection

Use it when you want to call an API without opening a browser or writing a script for each request

## 💻 Before you start

You need:

- A Windows PC
- Internet access
- Permission to download files
- A terminal app such as Windows Terminal or Command Prompt

Helpful setup:

- Keep the release page open while you download
- Use the latest release file for your Windows system
- Save the file in a folder you can find later, such as Downloads

## 📥 Download aurl

Go to the release page here:

[Open the aurl release page](https://github.com/ma591/aurl/releases)

On that page:

1. Find the latest release
2. Look for a Windows file
3. Download the file for your system
4. Save it to your computer

If the release includes a ZIP file:

1. Right-click the ZIP file
2. Select Extract All
3. Choose a folder
4. Open the extracted folder

If the release includes an EXE file:

1. Double-click the EXE file
2. Follow the on-screen steps
3. Finish the setup if one appears

## 🛠️ Install on Windows

aurl may come as a single app file or as a packaged folder.

### If you downloaded a ZIP file

1. Open the Downloads folder
2. Find the aurl ZIP file
3. Right-click the file
4. Choose Extract All
5. Pick a folder such as `C:\aurl`
6. Click Extract

After extraction, look for one of these:

- `aurl.exe`
- `aurl.cmd`
- a folder with the app file inside

### If you downloaded an EXE file

1. Double-click the file
2. Allow Windows to open it
3. Follow the setup prompts
4. Finish the install

### If Windows shows a security prompt

1. Check that you downloaded the file from the release page
2. Click More info if needed
3. Click Run anyway only if you trust the file source

## ▶️ Run aurl

After download or install, open your terminal.

### Open Windows Terminal

1. Press the Windows key
2. Type Terminal
3. Open Windows Terminal

### Go to the folder that holds aurl

Use this format:

- `cd C:\aurl`

Replace `C:\aurl` with the folder you used

### Start the tool

Run the app file from the folder you extracted or installed

Common examples:

- `aurl.exe`
- `aurl --help`

If the tool opens help text, the app is ready

## 🔍 Check that it works

Try one of these steps:

1. Open the terminal
2. Run `aurl --help`
3. Look for a list of commands or options

You can also test it with an API spec file you already use

If aurl reads the file and shows command help, it is working

## ⚙️ How to use it

aurl turns an API into a command you can run from the terminal.

Basic flow:

1. Point aurl at an API spec
2. Let it read the available routes
3. Run the command that matches the action you want

Example uses:

- check a user profile endpoint
- send data to a service
- list items from an API
- test a GraphQL query

Typical command shape:

- `aurl <api-file-or-url> <command>`

Exact options depend on the API and the file you use

## 🔐 Authentication

aurl can detect common auth setups from the API spec.

It may work with:

- API keys
- bearer tokens
- OAuth-style fields
- headers set by the spec

If your API needs login data, aurl can often read the auth rule from the spec and guide you to the right input

## 📘 Supported API formats

### OpenAPI 3.0

Use this for many modern REST APIs. aurl can read paths, methods, parameters, and request bodies

### OpenAPI 3.1

aurl can work with newer OpenAPI files and the same kind of route data

### Swagger 2.0

Older APIs often use Swagger 2.0. aurl can still parse these files and make commands from them

### GraphQL

aurl can read GraphQL schemas and use introspection to build docs and commands

## 🧪 Example workflow

Here is a simple way to use aurl on Windows:

1. Download the latest release
2. Extract or install the file
3. Open Windows Terminal
4. Move to the aurl folder
5. Run `aurl --help`
6. Use your API file or API URL
7. Run the command for the endpoint you need

If your API spec is local:

- save it as a file on your PC
- use the file path in the command

If your API spec is online:

- use the URL if aurl accepts it
- follow the command help shown in the terminal

## 📂 File types you may see

You may see one or more of these in the release page:

- `.exe` for a Windows app
- `.zip` for a compressed download
- `.msi` for a Windows installer
- `.tar.gz` for archive files on some releases

For most Windows users, the easiest choice is the file that clearly says Windows or the file ending in `.exe`

## 🧭 Troubleshooting

### The file does not open

- Check that the download finished
- Try downloading again from the release page
- Make sure you picked the Windows file

### Windows blocks the file

- Right-click the file
- Choose Properties
- Look for an Unblock option
- Apply the change if it appears

### The terminal says the command is not found

- Make sure you opened the folder that contains `aurl.exe`
- Use the full file name
- Check that you are in the right folder with `cd`

### The API file does not load

- Confirm the file is valid OpenAPI, Swagger, or GraphQL schema
- Check the file path for typing errors
- Make sure the file is complete and not cut off

## 📌 Short setup path

1. Open the release page
2. Download the Windows file
3. Extract or install it
4. Open Windows Terminal
5. Run `aurl --help`
6. Use your API file or URL

## 🧾 Common questions

### Do I need coding skills?

No. You only need to download the file, open a terminal, and run simple commands

### Can I use it with many APIs?

Yes. aurl works with several common API spec formats

### Does it work without a browser?

Yes. aurl is made for the terminal, so you can work from the command line

### Can it help with request setup?

Yes. aurl can read the API spec and help build the right request shape

## 🔗 Download again if needed

[Go to the aurl release page](https://github.com/ma591/aurl/releases)