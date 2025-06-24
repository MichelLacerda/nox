# 📦 Installation Guide for Nox

This guide explains how to build, install, and run the Nox interpreter on different platforms.

---

## 🪟 Windows

### 1. Add the local binary folder to your PATH (only once):

```sh
setx PATH "%USERPROFILE%\.local\bin;%PATH%"
```

### 2. Build for Windows:

```sh
go env -w GOOS=windows
make install-windows
```

---

## 🐧 Linux

### 1. Build for Linux:

```sh
go env -w GOOS=linux
make install-linux
```

---

## ▶️ Running Nox

To run the interpreter directly without installing:

```sh
make run
```

---

## 🧪 Testing Example Scripts

### On Windows:

```sh
make test-examples-windows
```

### On Linux:

```sh
make test-examples-linux
```

---

## ✅ Notes

- Ensure you have Go and Make installed.
- On Windows, run commands in **PowerShell** or **Command Prompt**.
- On Linux, ensure `make` has execution permission (`chmod +x` if needed).
