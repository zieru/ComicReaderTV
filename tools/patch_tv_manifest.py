#!/usr/bin/env python3
"""
patch_tv_manifest.py
Memodifikasi AndroidManifest.xml hasil ekstraksi APK agar 100% kompatibel dengan Android TV:
1. Menambahkan uses-feature leanback (required=false)
2. Menambahkan uses-feature touchscreen (required=false) agar TV tanpa layar sentuh tidak memblokir instalasi
3. Menambahkan category LEANBACK_LAUNCHER agar aplikasi muncul di Home Screen Android TV
4. Menambahkan atribut banner pada <application>
"""

import sys
import xml.etree.ElementTree as ET

def patch_manifest(manifest_path):
    ET.register_namespace('android', 'http://schemas.android.com/apk/res/android')
    tree = ET.parse(manifest_path)
    root = tree.getroot()
    android_ns = '{http://schemas.android.com/apk/res/android}'

    # Cek apakah leanback feature sudah ada
    has_leanback = any(elem.attrib.get(f'{android_ns}name') == 'android.software.leanback' for elem in root.findall('uses-feature'))
    if not has_leanback:
        f = ET.Element('uses-feature')
        f.set(f'{android_ns}name', 'android.software.leanback')
        f.set(f'{android_ns}required', 'false')
        root.insert(0, f)

    # Cek apakah touchscreen feature sudah ada
    has_touch = any(elem.attrib.get(f'{android_ns}name') == 'android.hardware.touchscreen' for elem in root.findall('uses-feature'))
    if not has_touch:
        f = ET.Element('uses-feature')
        f.set(f'{android_ns}name', 'android.hardware.touchscreen')
        f.set(f'{android_ns}required', 'false')
        root.insert(0, f)

    app = root.find('application')
    if app is not None:
        if f'{android_ns}banner' not in app.attrib and f'{android_ns}icon' in app.attrib:
            app.set(f'{android_ns}banner', app.attrib[f'{android_ns}icon'])
        
        act = app.find('activity')
        if act is not None:
            for ifilter in act.findall('intent-filter'):
                # Cek jika ada action MAIN
                has_main = any(action.attrib.get(f'{android_ns}name') == 'android.intent.action.MAIN' for action in ifilter.findall('action'))
                has_leanback_cat = any(cat.attrib.get(f'{android_ns}name') == 'android.intent.category.LEANBACK_LAUNCHER' for cat in ifilter.findall('category'))
                if has_main and not has_leanback_cat:
                    c = ET.Element('category')
                    c.set(f'{android_ns}name', 'android.intent.category.LEANBACK_LAUNCHER')
                    ifilter.append(c)

    tree.write(manifest_path, encoding='utf-8', xml_declaration=True)
    print(f"Successfully patched {manifest_path} for Android TV compatibility.")

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: patch_tv_manifest.py <path_to_AndroidManifest.xml>")
        sys.exit(1)
    patch_manifest(sys.argv[1])
