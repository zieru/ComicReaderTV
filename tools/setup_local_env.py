import os
import sys
import urllib.request
import zipfile
import shutil

SDK_DIR = r"C:\Users\Grapari_Infomedia\AppData\Local\Android\Sdk"
PLATFORMS_DIR = os.path.join(SDK_DIR, "platforms")
ANDROID_33_DIR = os.path.join(PLATFORMS_DIR, "android-33")
TOOLS_DIR = os.path.dirname(os.path.abspath(__file__))
APKTOOL_JAR = os.path.join(TOOLS_DIR, "apktool.jar")

def download_file(url, dest_path, desc=""):
    print(f"Downloading {desc or url}...")
    headers = {'User-Agent': 'Mozilla/5.0'}
    req = urllib.request.Request(url, headers=headers)
    with urllib.request.urlopen(req) as response, open(dest_path, 'wb') as out_file:
        total_size = int(response.info().get('Content-Length', 0))
        downloaded = 0
        chunk_size = 1024 * 1024  # 1MB
        while True:
            chunk = response.read(chunk_size)
            if not chunk:
                break
            out_file.write(chunk)
            downloaded += len(chunk)
            if total_size > 0:
                percent = (downloaded / total_size) * 100
                print(f"\r  {downloaded / (1024*1024):.1f}MB / {total_size / (1024*1024):.1f}MB ({percent:.1f}%)", end="")
            else:
                print(f"\r  {downloaded / (1024*1024):.1f}MB downloaded", end="")
    print("\nDownload complete.")

def setup_android_33():
    android_jar = os.path.join(ANDROID_33_DIR, "android.jar")
    if os.path.exists(android_jar):
        print(f"[OK] android-33 already exists at {ANDROID_33_DIR}")
        return

    os.makedirs(PLATFORMS_DIR, exist_ok=True)
    zip_path = os.path.join(SDK_DIR, "platform-33.zip")
    url = "https://dl.google.com/android/repository/platform-33_r02.zip"
    download_file(url, zip_path, "Android SDK Platform 33 (~64MB)")

    print("Extracting platform-33...")
    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
        zip_ref.extractall(PLATFORMS_DIR)
    
    extracted_folder = os.path.join(PLATFORMS_DIR, "android-13")
    if os.path.exists(extracted_folder):
        if os.path.exists(ANDROID_33_DIR):
            shutil.rmtree(ANDROID_33_DIR)
        os.rename(extracted_folder, ANDROID_33_DIR)
    
    if os.path.exists(zip_path):
        os.remove(zip_path)

    if os.path.exists(android_jar):
        print(f"[SUCCESS] android-33 installed at {ANDROID_33_DIR}")
    else:
        print(f"[WARNING] android.jar not found at {android_jar}")

def setup_apktool():
    if os.path.exists(APKTOOL_JAR) and os.path.getsize(APKTOOL_JAR) > 10 * 1024 * 1024:
        print(f"[OK] apktool.jar already exists at {APKTOOL_JAR}")
        return

    url = "https://github.com/iBotPeaches/Apktool/releases/download/v2.10.0/apktool_2.10.0.jar"
    download_file(url, APKTOOL_JAR, "Apktool v2.10.0 (~24MB)")
    if os.path.exists(APKTOOL_JAR):
        print(f"[SUCCESS] apktool.jar installed at {APKTOOL_JAR}")

def patch_d8():
    d8_bat = os.path.join(SDK_DIR, "build-tools", "30.0.3", "d8.bat")
    d8_jar = os.path.join(SDK_DIR, "build-tools", "30.0.3", "lib", "d8.jar")
    
    # Download modern R8/D8 (supports Java 21) if size is old (~4.8MB)
    if not os.path.exists(d8_jar) or os.path.getsize(d8_jar) < 10 * 1024 * 1024:
        r8_url = "https://dl.google.com/dl/android/maven2/com/android/tools/r8/8.2.33/r8-8.2.33.jar"
        download_file(r8_url, d8_jar, "Modern D8/R8 v8.2.33 (~15MB)")
    
    if os.path.exists(d8_bat):
        content = (
            "@echo off\r\n"
            "setlocal\r\n"
            "set \"jarpath=%~dp0lib\\d8.jar\"\r\n"
            "java -Xmx1024M -cp \"%jarpath%\" com.android.tools.r8.D8 %*\r\n"
        )
        with open(d8_bat, "w") as f:
            f.write(content)
        print(f"[SUCCESS] d8.bat patched for Java 21 compatibility at {d8_bat}")

if __name__ == "__main__":
    print("=== ComicReaderTV Local Build Environment Setup ===")
    setup_android_33()
    setup_apktool()
    patch_d8()
    print("=== Setup Finished ===")
