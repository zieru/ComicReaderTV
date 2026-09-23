.class public final Lorg/gioui/Gio;
.super Ljava/lang/Object;
.source "Gio.java"


# static fields
.field private static final handler:Landroid/os/Handler;

.field private static final initLock:Ljava/lang/Object;

.field private static jniLoaded:Z


# direct methods
.method static constructor <clinit>()V
    .locals 2

    .line 14
    new-instance v0, Ljava/lang/Object;

    invoke-direct {v0}, Ljava/lang/Object;-><init>()V

    sput-object v0, Lorg/gioui/Gio;->initLock:Ljava/lang/Object;

    .line 16
    new-instance v0, Landroid/os/Handler;

    invoke-static {}, Landroid/os/Looper;->getMainLooper()Landroid/os/Looper;

    move-result-object v1

    invoke-direct {v0, v1}, Landroid/os/Handler;-><init>(Landroid/os/Looper;)V

    sput-object v0, Lorg/gioui/Gio;->handler:Landroid/os/Handler;

    return-void
.end method

.method public constructor <init>()V
    .locals 0

    .line 13
    invoke-direct {p0}, Ljava/lang/Object;-><init>()V

    return-void
.end method

.method static synthetic access$000()V
    .locals 0

    .line 13
    invoke-static {}, Lorg/gioui/Gio;->scheduleMainFuncs()V

    return-void
.end method

.method public static declared-synchronized init(Landroid/content/Context;)V
    .locals 4

    const-class v0, Lorg/gioui/Gio;

    monitor-enter v0

    .line 26
    :try_start_0
    sget-object v1, Lorg/gioui/Gio;->initLock:Ljava/lang/Object;

    monitor-enter v1
    :try_end_0
    .catchall {:try_start_0 .. :try_end_0} :catchall_1

    .line 27
    :try_start_1
    sget-boolean v2, Lorg/gioui/Gio;->jniLoaded:Z

    if-eqz v2, :cond_0

    .line 28
    monitor-exit v1
    :try_end_1
    .catchall {:try_start_1 .. :try_end_1} :catchall_0

    monitor-exit v0

    return-void

    .line 30
    :cond_0
    :try_start_2
    invoke-virtual {p0}, Landroid/content/Context;->getFilesDir()Ljava/io/File;

    move-result-object v2

    invoke-virtual {v2}, Ljava/io/File;->getAbsolutePath()Ljava/lang/String;

    move-result-object v2
    :try_end_2
    .catchall {:try_start_2 .. :try_end_2} :catchall_0

    .line 33
    :try_start_3
    const-string v3, "UTF-8"

    invoke-virtual {v2, v3}, Ljava/lang/String;->getBytes(Ljava/lang/String;)[B

    move-result-object v2
    :try_end_3
    .catch Ljava/io/UnsupportedEncodingException; {:try_start_3 .. :try_end_3} :catch_0
    .catchall {:try_start_3 .. :try_end_3} :catchall_0

    .line 36
    nop

    .line 37
    :try_start_4
    const-string v3, "gio"

    invoke-static {v3}, Ljava/lang/System;->loadLibrary(Ljava/lang/String;)V

    .line 38
    invoke-static {v2, p0}, Lorg/gioui/Gio;->runGoMain([BLandroid/content/Context;)V

    .line 39
    const/4 p0, 0x1

    sput-boolean p0, Lorg/gioui/Gio;->jniLoaded:Z

    .line 40
    monitor-exit v1
    :try_end_4
    .catchall {:try_start_4 .. :try_end_4} :catchall_0

    .line 41
    monitor-exit v0

    return-void

    .line 34
    :catch_0
    move-exception p0

    .line 35
    :try_start_5
    new-instance v2, Ljava/lang/RuntimeException;

    invoke-direct {v2, p0}, Ljava/lang/RuntimeException;-><init>(Ljava/lang/Throwable;)V

    throw v2

    .line 40
    :catchall_0
    move-exception p0

    monitor-exit v1
    :try_end_5
    .catchall {:try_start_5 .. :try_end_5} :catchall_0

    :try_start_6
    throw p0
    :try_end_6
    .catchall {:try_start_6 .. :try_end_6} :catchall_1

    .line 25
    :catchall_1
    move-exception p0

    monitor-exit v0

    throw p0
.end method

.method static readClipboard(Landroid/content/Context;)Ljava/lang/String;
    .locals 3

    .line 51
    const-string v0, "clipboard"

    invoke-virtual {p0, v0}, Landroid/content/Context;->getSystemService(Ljava/lang/String;)Ljava/lang/Object;

    move-result-object v0

    check-cast v0, Landroid/content/ClipboardManager;

    .line 52
    invoke-virtual {v0}, Landroid/content/ClipboardManager;->getPrimaryClip()Landroid/content/ClipData;

    move-result-object v0

    .line 53
    if-eqz v0, :cond_1

    invoke-virtual {v0}, Landroid/content/ClipData;->getItemCount()I

    move-result v1

    const/4 v2, 0x1

    if-ge v1, v2, :cond_0

    goto :goto_0

    .line 56
    :cond_0
    const/4 v1, 0x0

    invoke-virtual {v0, v1}, Landroid/content/ClipData;->getItemAt(I)Landroid/content/ClipData$Item;

    move-result-object v0

    invoke-virtual {v0, p0}, Landroid/content/ClipData$Item;->coerceToText(Landroid/content/Context;)Ljava/lang/CharSequence;

    move-result-object p0

    invoke-interface {p0}, Ljava/lang/CharSequence;->toString()Ljava/lang/String;

    move-result-object p0

    return-object p0

    .line 54
    :cond_1
    :goto_0
    const/4 p0, 0x0

    return-object p0
.end method

.method private static native runGoMain([BLandroid/content/Context;)V
.end method

.method private static native scheduleMainFuncs()V
.end method

.method static wakeupMainThread()V
    .locals 2

    .line 60
    sget-object v0, Lorg/gioui/Gio;->handler:Landroid/os/Handler;

    new-instance v1, Lorg/gioui/Gio$1;

    invoke-direct {v1}, Lorg/gioui/Gio$1;-><init>()V

    invoke-virtual {v0, v1}, Landroid/os/Handler;->post(Ljava/lang/Runnable;)Z

    .line 65
    return-void
.end method

.method static writeClipboard(Landroid/content/Context;Ljava/lang/String;)V
    .locals 1

    .line 46
    const-string v0, "clipboard"

    invoke-virtual {p0, v0}, Landroid/content/Context;->getSystemService(Ljava/lang/String;)Ljava/lang/Object;

    move-result-object p0

    check-cast p0, Landroid/content/ClipboardManager;

    .line 47
    const/4 v0, 0x0

    invoke-static {v0, p1}, Landroid/content/ClipData;->newPlainText(Ljava/lang/CharSequence;Ljava/lang/CharSequence;)Landroid/content/ClipData;

    move-result-object p1

    invoke-virtual {p0, p1}, Landroid/content/ClipboardManager;->setPrimaryClip(Landroid/content/ClipData;)V

    .line 48
    return-void
.end method
