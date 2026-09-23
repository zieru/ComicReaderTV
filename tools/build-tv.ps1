# build-tv.ps1 - Automated 1-Click Local Build, Patch, Sign, & Install for ComicReaderTV
$ErrorActionPreference = "Stop"

$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$PROJECT_ROOT = Split-Path -Parent $SCRIPT_DIR
Set-Location $PROJECT_ROOT

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "   ComicReaderTV - Local Build & Deploy   " -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

# 1. Environment Variables
$env:ANDROID_HOME = "C:\Users\Grapari_Infomedia\AppData\Local\Android\Sdk"
$env:ANDROID_SDK_ROOT = $env:ANDROID_HOME
$env:ANDROID_NDK_ROOT = "D:\android-ndk-r30"
$env:ANDROID_NDK_HOME = $env:ANDROID_NDK_ROOT
$env:PATH = "$env:ANDROID_HOME\build-tools\30.0.3;$env:ANDROID_HOME\platform-tools;$env:PATH"

# 2. Check / Setup dependencies
Write-Host "`n[1/6] Verifying environment dependencies..." -ForegroundColor Yellow
python "$SCRIPT_DIR\setup_local_env.py"

# 3. Build Raw APK with gogio
Write-Host "`n[2/6] Building raw APK with gogio (target: android/arm)..." -ForegroundColor Yellow
$rawApk = "$PROJECT_ROOT\ComicReaderTV-raw.apk"
if (Test-Path $rawApk) { Remove-Item $rawApk -Force }

gogio -target android -arch arm -minsdk 21 -targetsdk 30 -name "ComicReaderTV" -appid com.comicreader.tv -icon "$PROJECT_ROOT\client\cmd\tv\appicon.png" -o $rawApk "$PROJECT_ROOT\client\cmd\tv"

if (-not (Test-Path $rawApk)) {
    Write-Error "Failed: ComicReaderTV-raw.apk was not created by gogio!"
}
Write-Host "[OK] Raw APK created successfully." -ForegroundColor Green

# 4. Decompile with Apktool
Write-Host "`n[3/6] Decompiling APK with Apktool..." -ForegroundColor Yellow
$unpackedDir = "$PROJECT_ROOT\apk_unpacked"
if (Test-Path $unpackedDir) { Remove-Item $unpackedDir -Recurse -Force }

java -jar "$SCRIPT_DIR\apktool.jar" d -f $rawApk -o $unpackedDir -q

# 5. Patch AndroidManifest for Android TV (Leanback + Cleartext HTTP + Permissions)
Write-Host "`n[4/6] Patching AndroidManifest.xml for Android TV..." -ForegroundColor Yellow
python "$SCRIPT_DIR\patch_tv_manifest.py" "$unpackedDir\AndroidManifest.xml"

# 6. Rebuild APK
Write-Host "`n[5/6] Rebuilding modified APK..." -ForegroundColor Yellow
$unalignedApk = "$PROJECT_ROOT\ComicReaderTV-unaligned.apk"
if (Test-Path $unalignedApk) { Remove-Item $unalignedApk -Force }

java -jar "$SCRIPT_DIR\apktool.jar" b $unpackedDir -o $unalignedApk -q

# 7. Zipalign & Sign
Write-Host "`n[6/6] Zipaligning & Signing APK..." -ForegroundColor Yellow
$alignedApk = "$PROJECT_ROOT\ComicReaderTV-aligned.apk"
$finalApk = "$PROJECT_ROOT\ComicReaderTV-debug.apk"
if (Test-Path $alignedApk) { Remove-Item $alignedApk -Force }
if (Test-Path $finalApk) { Remove-Item $finalApk -Force }

& "$env:ANDROID_HOME\build-tools\30.0.3\zipalign.exe" -p -f 4 $unalignedApk $alignedApk
& "$env:ANDROID_HOME\build-tools\30.0.3\apksigner.bat" sign --ks "$SCRIPT_DIR\release.keystore" --ks-pass pass:comicreadertv --out $finalApk $alignedApk

Write-Host "`n=======================================================" -ForegroundColor Green
Write-Host "   BUILD SUCCESS: ComicReaderTV-debug.apk is ready!   " -ForegroundColor Green
Write-Host "=======================================================" -ForegroundColor Green

# 8. Check ADB device and auto-install
$adbDevices = adb devices
$targetDevice = "10.0.1.5:5555"
if ($adbDevices -match $targetDevice) {
    Write-Host "`n[DEPLOY] Installing to TV ($targetDevice)..." -ForegroundColor Cyan
    adb -s $targetDevice install -r $finalApk
    Write-Host "[DEPLOY] Launching ComicReaderTV on TV..." -ForegroundColor Cyan
    adb -s $targetDevice shell am start -n com.comicreader.tv/org.gioui.GioActivity
    Write-Host "[DEPLOY] Successfully launched on TV!" -ForegroundColor Green
} else {
    Write-Host "`n[INFO] Target device $targetDevice not detected. APK is ready at: $finalApk" -ForegroundColor Yellow
}
