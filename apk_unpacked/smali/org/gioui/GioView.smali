.class public final Lorg/gioui/GioView;
.super Landroid/view/SurfaceView;
.source "GioView.java"

# interfaces
.implements Landroid/view/Choreographer$FrameCallback;


# annotations
.annotation system Ldalvik/annotation/MemberClasses;
    value = {
        Lorg/gioui/GioView$Snippet;,
        Lorg/gioui/GioView$Bar;,
        Lorg/gioui/GioView$GioInputConnection;
    }
.end annotation


# static fields
.field private static jniLoaded:Z


# instance fields
.field private final accessManager:Landroid/view/accessibility/AccessibilityManager;

.field private final imm:Landroid/view/inputmethod/InputMethodManager;

.field private keyboardHint:I

.field private nhandle:J

.field private final scrollXScale:F

.field private final scrollYScale:F

.field private final surfCallbacks:Landroid/view/SurfaceHolder$Callback;


# direct methods
.method public constructor <init>(Landroid/content/Context;)V
    .locals 1

    .line 74
    const/4 v0, 0x0

    invoke-direct {p0, p1, v0}, Lorg/gioui/GioView;-><init>(Landroid/content/Context;Landroid/util/AttributeSet;)V

    .line 75
    return-void
.end method

.method public constructor <init>(Landroid/content/Context;Landroid/util/AttributeSet;)V
    .locals 4

    .line 78
    invoke-direct {p0, p1, p2}, Landroid/view/SurfaceView;-><init>(Landroid/content/Context;Landroid/util/AttributeSet;)V

    .line 79
    nop

    .line 80
    const/16 p2, 0x300

    invoke-virtual {p0, p2}, Lorg/gioui/GioView;->setSystemUiVisibility(I)V

    .line 82
    new-instance p2, Landroid/view/WindowManager$LayoutParams;

    const/4 v0, -0x1

    invoke-direct {p2, v0, v0}, Landroid/view/WindowManager$LayoutParams;-><init>(II)V

    invoke-virtual {p0, p2}, Lorg/gioui/GioView;->setLayoutParams(Landroid/view/ViewGroup$LayoutParams;)V

    .line 85
    invoke-virtual {p1}, Landroid/content/Context;->getApplicationContext()Landroid/content/Context;

    move-result-object p2

    invoke-static {p2}, Lorg/gioui/Gio;->init(Landroid/content/Context;)V

    .line 89
    const/4 p2, 0x0

    invoke-static {p2, p2, p2, p2}, Landroid/graphics/Color;->argb(IIII)I

    move-result v0

    invoke-virtual {p0, v0}, Lorg/gioui/GioView;->setBackgroundColor(I)V

    .line 91
    invoke-static {p1}, Landroid/view/ViewConfiguration;->get(Landroid/content/Context;)Landroid/view/ViewConfiguration;

    move-result-object v0

    .line 92
    sget v1, Landroid/os/Build$VERSION;->SDK_INT:I

    const/16 v2, 0x1a

    const/4 v3, 0x1

    if-lt v1, v2, :cond_0

    .line 93
    invoke-virtual {v0}, Landroid/view/ViewConfiguration;->getScaledHorizontalScrollFactor()F

    move-result v1

    iput v1, p0, Lorg/gioui/GioView;->scrollXScale:F

    .line 94
    invoke-virtual {v0}, Landroid/view/ViewConfiguration;->getScaledVerticalScrollFactor()F

    move-result v0

    iput v0, p0, Lorg/gioui/GioView;->scrollYScale:F

    .line 97
    invoke-virtual {p0, p2}, Lorg/gioui/GioView;->setDefaultFocusHighlightEnabled(Z)V

    goto :goto_0

    .line 99
    :cond_0
    nop

    .line 100
    nop

    .line 103
    invoke-virtual {p0}, Lorg/gioui/GioView;->getResources()Landroid/content/res/Resources;

    move-result-object p2

    invoke-virtual {p2}, Landroid/content/res/Resources;->getDisplayMetrics()Landroid/util/DisplayMetrics;

    move-result-object p2

    .line 100
    const/high16 v0, 0x42400000    # 48.0f

    invoke-static {v3, v0, p2}, Landroid/util/TypedValue;->applyDimension(IFLandroid/util/DisplayMetrics;)F

    move-result p2

    .line 105
    iput p2, p0, Lorg/gioui/GioView;->scrollXScale:F

    .line 106
    iput p2, p0, Lorg/gioui/GioView;->scrollYScale:F

    .line 109
    :goto_0
    invoke-direct {p0}, Lorg/gioui/GioView;->setHighRefreshRate()V

    .line 111
    const-string p2, "accessibility"

    invoke-virtual {p1, p2}, Landroid/content/Context;->getSystemService(Ljava/lang/String;)Ljava/lang/Object;

    move-result-object p2

    check-cast p2, Landroid/view/accessibility/AccessibilityManager;

    iput-object p2, p0, Lorg/gioui/GioView;->accessManager:Landroid/view/accessibility/AccessibilityManager;

    .line 112
    const-string p2, "input_method"

    invoke-virtual {p1, p2}, Landroid/content/Context;->getSystemService(Ljava/lang/String;)Ljava/lang/Object;

    move-result-object p1

    check-cast p1, Landroid/view/inputmethod/InputMethodManager;

    iput-object p1, p0, Lorg/gioui/GioView;->imm:Landroid/view/inputmethod/InputMethodManager;

    .line 113
    invoke-static {p0}, Lorg/gioui/GioView;->onCreateView(Lorg/gioui/GioView;)J

    move-result-wide p1

    iput-wide p1, p0, Lorg/gioui/GioView;->nhandle:J

    .line 114
    invoke-virtual {p0, v3}, Lorg/gioui/GioView;->setFocusable(Z)V

    .line 115
    invoke-virtual {p0, v3}, Lorg/gioui/GioView;->setFocusableInTouchMode(Z)V

    .line 116
    new-instance p1, Lorg/gioui/GioView$1;

    invoke-direct {p1, p0}, Lorg/gioui/GioView$1;-><init>(Lorg/gioui/GioView;)V

    iput-object p1, p0, Lorg/gioui/GioView;->surfCallbacks:Landroid/view/SurfaceHolder$Callback;

    .line 127
    invoke-virtual {p0}, Lorg/gioui/GioView;->getHolder()Landroid/view/SurfaceHolder;

    move-result-object p2

    invoke-interface {p2, p1}, Landroid/view/SurfaceHolder;->addCallback(Landroid/view/SurfaceHolder$Callback;)V

    .line 128
    return-void
.end method

.method static synthetic access$000(Lorg/gioui/GioView;)J
    .locals 2

    .line 61
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    return-wide v0
.end method

.method static synthetic access$100(JLandroid/view/Surface;)V
    .locals 0

    .line 61
    invoke-static {p0, p1, p2}, Lorg/gioui/GioView;->onSurfaceChanged(JLandroid/view/Surface;)V

    return-void
.end method

.method static synthetic access$1000(JII)V
    .locals 0

    .line 61
    invoke-static {p0, p1, p2, p3}, Lorg/gioui/GioView;->imeSetSnippet(JII)V

    return-void
.end method

.method static synthetic access$1100(JIIZJ)V
    .locals 0

    .line 61
    invoke-static/range {p0 .. p6}, Lorg/gioui/GioView;->onKeyEvent(JIIZJ)V

    return-void
.end method

.method static synthetic access$1200(J)I
    .locals 0

    .line 61
    invoke-static {p0, p1}, Lorg/gioui/GioView;->imeComposingStart(J)I

    move-result p0

    return p0
.end method

.method static synthetic access$1300(J)I
    .locals 0

    .line 61
    invoke-static {p0, p1}, Lorg/gioui/GioView;->imeComposingEnd(J)I

    move-result p0

    return p0
.end method

.method static synthetic access$1400(JIILjava/lang/String;)I
    .locals 0

    .line 61
    invoke-static {p0, p1, p2, p3, p4}, Lorg/gioui/GioView;->imeReplace(JIILjava/lang/String;)I

    move-result p0

    return p0
.end method

.method static synthetic access$1500(JII)I
    .locals 0

    .line 61
    invoke-static {p0, p1, p2, p3}, Lorg/gioui/GioView;->imeSetSelection(JII)I

    move-result p0

    return p0
.end method

.method static synthetic access$1700(JIIILandroid/view/accessibility/AccessibilityNodeInfo;)Landroid/view/accessibility/AccessibilityNodeInfo;
    .locals 0

    .line 61
    invoke-static/range {p0 .. p5}, Lorg/gioui/GioView;->initializeAccessibilityNodeInfo(JIIILandroid/view/accessibility/AccessibilityNodeInfo;)Landroid/view/accessibility/AccessibilityNodeInfo;

    move-result-object p0

    return-object p0
.end method

.method static synthetic access$1800(JI)V
    .locals 0

    .line 61
    invoke-static {p0, p1, p2}, Lorg/gioui/GioView;->onA11yFocus(JI)V

    return-void
.end method

.method static synthetic access$1900(JI)V
    .locals 0

    .line 61
    invoke-static {p0, p1, p2}, Lorg/gioui/GioView;->onClearA11yFocus(JI)V

    return-void
.end method

.method static synthetic access$200(J)V
    .locals 0

    .line 61
    invoke-static {p0, p1}, Lorg/gioui/GioView;->onSurfaceDestroyed(J)V

    return-void
.end method

.method static synthetic access$400(J)I
    .locals 0

    .line 61
    invoke-static {p0, p1}, Lorg/gioui/GioView;->imeSelectionStart(J)I

    move-result p0

    return p0
.end method

.method static synthetic access$500(J)I
    .locals 0

    .line 61
    invoke-static {p0, p1}, Lorg/gioui/GioView;->imeSelectionEnd(J)I

    move-result p0

    return p0
.end method

.method static synthetic access$600(JI)I
    .locals 0

    .line 61
    invoke-static {p0, p1, p2}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result p0

    return p0
.end method

.method static synthetic access$700(JI)I
    .locals 0

    .line 61
    invoke-static {p0, p1, p2}, Lorg/gioui/GioView;->imeToRunes(JI)I

    move-result p0

    return p0
.end method

.method static synthetic access$800(JII)I
    .locals 0

    .line 61
    invoke-static {p0, p1, p2, p3}, Lorg/gioui/GioView;->imeSetComposingRegion(JII)I

    move-result p0

    return p0
.end method

.method static synthetic access$900(Lorg/gioui/GioView;)Lorg/gioui/GioView$Snippet;
    .locals 0

    .line 61
    invoke-direct {p0}, Lorg/gioui/GioView;->getSnippet()Lorg/gioui/GioView$Snippet;

    move-result-object p0

    return-object p0
.end method

.method private dispatchMotionEvent(Landroid/view/MotionEvent;)V
    .locals 21

    .line 369
    move-object/from16 v0, p0

    move-object/from16 v1, p1

    iget-wide v2, v0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v4, 0x0

    cmp-long v6, v2, v4

    if-nez v6, :cond_0

    .line 370
    return-void

    .line 372
    :cond_0
    const/4 v2, 0x0

    const/4 v3, 0x0

    :goto_0
    invoke-virtual/range {p1 .. p1}, Landroid/view/MotionEvent;->getHistorySize()I

    move-result v4

    const/16 v5, 0x9

    const/16 v6, 0xa

    if-ge v3, v4, :cond_2

    .line 373
    invoke-virtual {v1, v3}, Landroid/view/MotionEvent;->getHistoricalEventTime(I)J

    move-result-wide v19

    .line 374
    const/4 v4, 0x0

    :goto_1
    invoke-virtual/range {p1 .. p1}, Landroid/view/MotionEvent;->getPointerCount()I

    move-result v7

    if-ge v4, v7, :cond_1

    .line 375
    iget-wide v7, v0, Lorg/gioui/GioView;->nhandle:J

    const/4 v9, 0x2

    .line 378
    invoke-virtual {v1, v4}, Landroid/view/MotionEvent;->getPointerId(I)I

    move-result v10

    .line 379
    invoke-virtual {v1, v4}, Landroid/view/MotionEvent;->getToolType(I)I

    move-result v11

    .line 380
    invoke-virtual {v1, v4, v3}, Landroid/view/MotionEvent;->getHistoricalX(II)F

    move-result v12

    .line 381
    invoke-virtual {v1, v4, v3}, Landroid/view/MotionEvent;->getHistoricalY(II)F

    move-result v13

    iget v14, v0, Lorg/gioui/GioView;->scrollXScale:F

    .line 382
    invoke-virtual {v1, v6, v4, v3}, Landroid/view/MotionEvent;->getHistoricalAxisValue(III)F

    move-result v15

    mul-float v14, v14, v15

    iget v15, v0, Lorg/gioui/GioView;->scrollYScale:F

    .line 383
    invoke-virtual {v1, v5, v4, v3}, Landroid/view/MotionEvent;->getHistoricalAxisValue(III)F

    move-result v16

    mul-float v15, v15, v16

    .line 384
    invoke-virtual/range {p1 .. p1}, Landroid/view/MotionEvent;->getButtonState()I

    move-result v16

    .line 375
    move-wide/from16 v17, v19

    invoke-static/range {v7 .. v18}, Lorg/gioui/GioView;->onTouchEvent(JIIIFFFFIJ)V

    .line 374
    add-int/lit8 v4, v4, 0x1

    goto :goto_1

    .line 372
    :cond_1
    add-int/lit8 v3, v3, 0x1

    goto :goto_0

    .line 388
    :cond_2
    invoke-virtual/range {p1 .. p1}, Landroid/view/MotionEvent;->getActionMasked()I

    move-result v3

    .line 389
    invoke-virtual/range {p1 .. p1}, Landroid/view/MotionEvent;->getActionIndex()I

    move-result v4

    .line 390
    nop

    :goto_2
    invoke-virtual/range {p1 .. p1}, Landroid/view/MotionEvent;->getPointerCount()I

    move-result v7

    if-ge v2, v7, :cond_4

    .line 391
    nop

    .line 392
    if-ne v2, v4, :cond_3

    .line 393
    move v10, v3

    goto :goto_3

    .line 392
    :cond_3
    const/4 v7, 0x2

    const/4 v10, 0x2

    .line 395
    :goto_3
    iget-wide v8, v0, Lorg/gioui/GioView;->nhandle:J

    .line 398
    invoke-virtual {v1, v2}, Landroid/view/MotionEvent;->getPointerId(I)I

    move-result v11

    .line 399
    invoke-virtual {v1, v2}, Landroid/view/MotionEvent;->getToolType(I)I

    move-result v12

    .line 400
    invoke-virtual {v1, v2}, Landroid/view/MotionEvent;->getX(I)F

    move-result v13

    invoke-virtual {v1, v2}, Landroid/view/MotionEvent;->getY(I)F

    move-result v14

    iget v7, v0, Lorg/gioui/GioView;->scrollXScale:F

    .line 401
    invoke-virtual {v1, v6, v2}, Landroid/view/MotionEvent;->getAxisValue(II)F

    move-result v15

    mul-float v15, v15, v7

    iget v7, v0, Lorg/gioui/GioView;->scrollYScale:F

    .line 402
    invoke-virtual {v1, v5, v2}, Landroid/view/MotionEvent;->getAxisValue(II)F

    move-result v16

    mul-float v16, v16, v7

    .line 403
    invoke-virtual/range {p1 .. p1}, Landroid/view/MotionEvent;->getButtonState()I

    move-result v17

    .line 404
    invoke-virtual/range {p1 .. p1}, Landroid/view/MotionEvent;->getEventTime()J

    move-result-wide v18

    .line 395
    invoke-static/range {v8 .. v19}, Lorg/gioui/GioView;->onTouchEvent(JIIIFFFFIJ)V

    .line 390
    add-int/lit8 v2, v2, 0x1

    goto :goto_2

    .line 406
    :cond_4
    return-void
.end method

.method private getSnippet()Lorg/gioui/GioView$Snippet;
    .locals 3

    .line 782
    new-instance v0, Lorg/gioui/GioView$Snippet;

    const/4 v1, 0x0

    invoke-direct {v0, v1}, Lorg/gioui/GioView$Snippet;-><init>(Lorg/gioui/GioView$1;)V

    .line 783
    iget-wide v1, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v1, v2}, Lorg/gioui/GioView;->imeSnippet(J)Ljava/lang/String;

    move-result-object v1

    iput-object v1, v0, Lorg/gioui/GioView$Snippet;->snippet:Ljava/lang/String;

    .line 784
    iget-wide v1, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v1, v2}, Lorg/gioui/GioView;->imeSnippetStart(J)I

    move-result v1

    iput v1, v0, Lorg/gioui/GioView$Snippet;->offset:I

    .line 785
    return-object v0
.end method

.method private static native imeComposingEnd(J)I
.end method

.method private static native imeComposingStart(J)I
.end method

.method private static native imeReplace(JIILjava/lang/String;)I
.end method

.method private static native imeSelectionEnd(J)I
.end method

.method private static native imeSelectionStart(J)I
.end method

.method private static native imeSetComposingRegion(JII)I
.end method

.method private static native imeSetSelection(JII)I
.end method

.method private static native imeSetSnippet(JII)V
.end method

.method private static native imeSnippet(J)Ljava/lang/String;
.end method

.method private static native imeSnippetStart(J)I
.end method

.method private static native imeToRunes(JI)I
.end method

.method private static native imeToUTF16(JI)I
.end method

.method private static native initializeAccessibilityNodeInfo(JIIILandroid/view/accessibility/AccessibilityNodeInfo;)Landroid/view/accessibility/AccessibilityNodeInfo;
.end method

.method private static native onA11yFocus(JI)V
.end method

.method private static native onBack(J)Z
.end method

.method private static native onClearA11yFocus(JI)V
.end method

.method private static native onConfigurationChanged(J)V
.end method

.method private static native onCreateView(Lorg/gioui/GioView;)J
.end method

.method private static native onDestroyView(J)V
.end method

.method private static native onExitTouchExploration(J)V
.end method

.method private static native onFocusChange(JZ)V
.end method

.method private static native onFrameCallback(J)V
.end method

.method private static native onKeyEvent(JIIZJ)V
.end method

.method public static native onLowMemory()V
.end method

.method private static native onOpenURI(JLjava/lang/String;)V
.end method

.method private static native onStartView(J)V
.end method

.method private static native onStopView(J)V
.end method

.method private static native onSurfaceChanged(JLandroid/view/Surface;)V
.end method

.method private static native onSurfaceDestroyed(J)V
.end method

.method private static native onTouchEvent(JIIIFFFFIJ)V
.end method

.method private static native onTouchExploration(JFF)V
.end method

.method private static native onWindowInsets(JIIII)V
.end method

.method private setBarColor(Lorg/gioui/GioView$Bar;II)V
    .locals 4

    .line 197
    nop

    .line 201
    invoke-virtual {p0}, Lorg/gioui/GioView;->getContext()Landroid/content/Context;

    move-result-object v0

    check-cast v0, Landroid/app/Activity;

    invoke-virtual {v0}, Landroid/app/Activity;->getWindow()Landroid/view/Window;

    move-result-object v0

    .line 206
    invoke-virtual {p1}, Lorg/gioui/GioView$Bar;->ordinal()I

    move-result p1

    packed-switch p1, :pswitch_data_0

    .line 218
    new-instance p1, Ljava/lang/RuntimeException;

    const-string p2, "invalid bar type"

    invoke-direct {p1, p2}, Ljava/lang/RuntimeException;-><init>(Ljava/lang/String;)V

    throw p1

    .line 208
    :pswitch_0
    nop

    .line 209
    nop

    .line 210
    invoke-virtual {v0, p2}, Landroid/view/Window;->setStatusBarColor(I)V

    .line 211
    const/16 p1, 0x8

    const/16 p2, 0x2000

    goto :goto_0

    .line 213
    :pswitch_1
    nop

    .line 214
    nop

    .line 215
    invoke-virtual {v0, p2}, Landroid/view/Window;->setNavigationBarColor(I)V

    .line 216
    const/16 p1, 0x10

    const/16 p2, 0x10

    .line 221
    :goto_0
    sget v1, Landroid/os/Build$VERSION;->SDK_INT:I

    const/16 v2, 0x17

    if-ge v1, v2, :cond_0

    .line 222
    return-void

    .line 225
    :cond_0
    sget v1, Landroid/os/Build$VERSION;->SDK_INT:I

    const/16 v2, 0x1e

    const/16 v3, 0x80

    if-ge v1, v2, :cond_2

    .line 226
    invoke-virtual {p0}, Lorg/gioui/GioView;->getSystemUiVisibility()I

    move-result p1

    .line 227
    if-le p3, v3, :cond_1

    .line 228
    or-int/2addr p1, p2

    goto :goto_1

    .line 230
    :cond_1
    not-int p2, p2

    and-int/2addr p1, p2

    .line 232
    :goto_1
    invoke-virtual {p0, p1}, Lorg/gioui/GioView;->setSystemUiVisibility(I)V

    .line 233
    return-void

    .line 236
    :cond_2
    invoke-virtual {v0}, Landroid/view/Window;->getInsetsController()Landroid/view/WindowInsetsController;

    move-result-object p2

    .line 237
    if-nez p2, :cond_3

    .line 238
    return-void

    .line 240
    :cond_3
    if-le p3, v3, :cond_4

    .line 241
    invoke-interface {p2, p1, p1}, Landroid/view/WindowInsetsController;->setSystemBarsAppearance(II)V

    goto :goto_2

    .line 243
    :cond_4
    const/4 p3, 0x0

    invoke-interface {p2, p3, p1}, Landroid/view/WindowInsetsController;->setSystemBarsAppearance(II)V

    .line 245
    :goto_2
    return-void

    :pswitch_data_0
    .packed-switch 0x0
        :pswitch_1
        :pswitch_0
    .end packed-switch
.end method

.method private setCursor(I)V
    .locals 2

    .line 161
    sget v0, Landroid/os/Build$VERSION;->SDK_INT:I

    const/16 v1, 0x18

    if-ge v0, v1, :cond_0

    .line 162
    return-void

    .line 164
    :cond_0
    invoke-virtual {p0}, Lorg/gioui/GioView;->getContext()Landroid/content/Context;

    move-result-object v0

    invoke-static {v0, p1}, Landroid/view/PointerIcon;->getSystemIcon(Landroid/content/Context;I)Landroid/view/PointerIcon;

    move-result-object p1

    .line 165
    invoke-virtual {p0, p1}, Lorg/gioui/GioView;->setPointerIcon(Landroid/view/PointerIcon;)V

    .line 166
    return-void
.end method

.method private setFullscreen(Z)V
    .locals 1

    .line 176
    invoke-virtual {p0}, Lorg/gioui/GioView;->getSystemUiVisibility()I

    move-result v0

    .line 177
    if-eqz p1, :cond_0

    .line 178
    or-int/lit16 p1, v0, 0x1000

    .line 179
    or-int/lit8 p1, p1, 0x2

    .line 180
    or-int/lit8 p1, p1, 0x4

    .line 181
    or-int/lit16 p1, p1, 0x400

    goto :goto_0

    .line 183
    :cond_0
    and-int/lit16 p1, v0, -0x1001

    .line 184
    and-int/lit8 p1, p1, -0x3

    .line 185
    and-int/lit8 p1, p1, -0x5

    .line 186
    and-int/lit16 p1, p1, -0x401

    .line 188
    :goto_0
    invoke-virtual {p0, p1}, Lorg/gioui/GioView;->setSystemUiVisibility(I)V

    .line 189
    return-void
.end method

.method private setHighRefreshRate()V
    .locals 17

    .line 256
    sget v0, Landroid/os/Build$VERSION;->SDK_INT:I

    const/16 v1, 0x1e

    if-ge v0, v1, :cond_0

    .line 257
    return-void

    .line 260
    :cond_0
    invoke-virtual/range {p0 .. p0}, Lorg/gioui/GioView;->getContext()Landroid/content/Context;

    move-result-object v0

    .line 261
    invoke-virtual {v0}, Landroid/content/Context;->getDisplay()Landroid/view/Display;

    move-result-object v1

    .line 262
    invoke-virtual {v1}, Landroid/view/Display;->getSupportedModes()[Landroid/view/Display$Mode;

    move-result-object v2

    .line 263
    array-length v3, v2

    const/4 v4, 0x1

    if-gt v3, v4, :cond_1

    .line 265
    return-void

    .line 268
    :cond_1
    invoke-virtual {v1}, Landroid/view/Display;->getMode()Landroid/view/Display$Mode;

    move-result-object v1

    .line 269
    invoke-virtual {v1}, Landroid/view/Display$Mode;->getPhysicalWidth()I

    move-result v3

    .line 270
    invoke-virtual {v1}, Landroid/view/Display$Mode;->getPhysicalHeight()I

    move-result v1

    .line 272
    nop

    .line 273
    nop

    .line 274
    nop

    .line 275
    nop

    .line 276
    array-length v5, v2

    const/4 v7, -0x1

    const/high16 v8, -0x40800000    # -1.0f

    const/4 v9, 0x0

    const/4 v10, -0x1

    const/high16 v11, -0x40800000    # -1.0f

    const/high16 v12, -0x40800000    # -1.0f

    const/high16 v13, -0x40800000    # -1.0f

    :goto_0
    if-ge v9, v5, :cond_9

    aget-object v14, v2, v9

    .line 277
    invoke-virtual {v14}, Landroid/view/Display$Mode;->getRefreshRate()F

    move-result v15

    .line 278
    invoke-virtual {v14}, Landroid/view/Display$Mode;->getPhysicalWidth()I

    move-result v4

    int-to-float v4, v4

    .line 279
    invoke-virtual {v14}, Landroid/view/Display$Mode;->getPhysicalHeight()I

    move-result v6

    int-to-float v6, v6

    .line 281
    cmpl-float v16, v11, v8

    if-eqz v16, :cond_2

    cmpg-float v16, v15, v11

    if-gez v16, :cond_3

    .line 282
    :cond_2
    move v11, v15

    .line 284
    :cond_3
    cmpl-float v16, v12, v8

    if-eqz v16, :cond_4

    cmpl-float v16, v15, v12

    if-lez v16, :cond_5

    .line 285
    :cond_4
    move v12, v15

    .line 288
    :cond_5
    cmpl-float v16, v13, v8

    if-eqz v16, :cond_7

    cmpl-float v16, v15, v13

    if-lez v16, :cond_6

    goto :goto_1

    :cond_6
    const/16 v16, 0x0

    goto :goto_2

    :cond_7
    :goto_1
    const/16 v16, 0x1

    .line 289
    :goto_2
    int-to-float v8, v3

    cmpl-float v4, v4, v8

    if-nez v4, :cond_8

    int-to-float v4, v1

    cmpl-float v4, v6, v4

    if-nez v4, :cond_8

    if-eqz v16, :cond_8

    .line 290
    invoke-virtual {v14}, Landroid/view/Display$Mode;->getModeId()I

    move-result v10

    .line 291
    nop

    .line 292
    move v13, v15

    .line 276
    :cond_8
    add-int/lit8 v9, v9, 0x1

    const/4 v4, 0x1

    const/high16 v8, -0x40800000    # -1.0f

    goto :goto_0

    .line 296
    :cond_9
    if-ne v10, v7, :cond_a

    .line 298
    return-void

    .line 301
    :cond_a
    cmpl-float v1, v11, v12

    if-nez v1, :cond_b

    .line 303
    return-void

    .line 306
    :cond_b
    check-cast v0, Landroid/app/Activity;

    invoke-virtual {v0}, Landroid/app/Activity;->getWindow()Landroid/view/Window;

    move-result-object v0

    .line 307
    invoke-virtual {v0}, Landroid/view/Window;->getAttributes()Landroid/view/WindowManager$LayoutParams;

    move-result-object v1

    .line 308
    iput v10, v1, Landroid/view/WindowManager$LayoutParams;->preferredDisplayModeId:I

    .line 309
    invoke-virtual {v0, v1}, Landroid/view/Window;->setAttributes(Landroid/view/WindowManager$LayoutParams;)V

    .line 310
    return-void
.end method

.method private setNavigationColor(II)V
    .locals 1

    .line 252
    sget-object v0, Lorg/gioui/GioView$Bar;->NAVIGATION:Lorg/gioui/GioView$Bar;

    invoke-direct {p0, v0, p1, p2}, Lorg/gioui/GioView;->setBarColor(Lorg/gioui/GioView$Bar;II)V

    .line 253
    return-void
.end method

.method private setOrientation(II)V
    .locals 0

    .line 169
    nop

    .line 172
    invoke-virtual {p0}, Lorg/gioui/GioView;->getContext()Landroid/content/Context;

    move-result-object p2

    check-cast p2, Landroid/app/Activity;

    invoke-virtual {p2, p1}, Landroid/app/Activity;->setRequestedOrientation(I)V

    .line 173
    return-void
.end method

.method private setStatusColor(II)V
    .locals 1

    .line 248
    sget-object v0, Lorg/gioui/GioView$Bar;->STATUS:Lorg/gioui/GioView$Bar;

    invoke-direct {p0, v0, p1, p2}, Lorg/gioui/GioView;->setBarColor(Lorg/gioui/GioView$Bar;II)V

    .line 249
    return-void
.end method


# virtual methods
.method public backPressed()Z
    .locals 5

    .line 509
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-nez v4, :cond_0

    .line 510
    const/4 v0, 0x0

    return v0

    .line 512
    :cond_0
    invoke-static {v0, v1}, Lorg/gioui/GioView;->onBack(J)Z

    move-result v0

    return v0
.end method

.method public configurationChanged()V
    .locals 5

    .line 503
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-eqz v4, :cond_0

    .line 504
    invoke-static {v0, v1}, Lorg/gioui/GioView;->onConfigurationChanged(J)V

    .line 506
    :cond_0
    return-void
.end method

.method public destroy()V
    .locals 5

    .line 491
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-eqz v4, :cond_0

    .line 492
    invoke-static {v0, v1}, Lorg/gioui/GioView;->onDestroyView(J)V

    .line 494
    :cond_0
    return-void
.end method

.method protected dispatchHoverEvent(Landroid/view/MotionEvent;)Z
    .locals 3

    .line 322
    iget-object v0, p0, Lorg/gioui/GioView;->accessManager:Landroid/view/accessibility/AccessibilityManager;

    invoke-virtual {v0}, Landroid/view/accessibility/AccessibilityManager;->isTouchExplorationEnabled()Z

    move-result v0

    if-nez v0, :cond_0

    .line 323
    invoke-super {p0, p1}, Landroid/view/SurfaceView;->dispatchHoverEvent(Landroid/view/MotionEvent;)Z

    move-result p1

    return p1

    .line 325
    :cond_0
    invoke-virtual {p1}, Landroid/view/MotionEvent;->getAction()I

    move-result v0

    packed-switch v0, :pswitch_data_0

    :pswitch_0
    goto :goto_0

    .line 332
    :pswitch_1
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v0, v1}, Lorg/gioui/GioView;->onExitTouchExploration(J)V

    goto :goto_0

    .line 329
    :pswitch_2
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-virtual {p1}, Landroid/view/MotionEvent;->getX()F

    move-result v2

    invoke-virtual {p1}, Landroid/view/MotionEvent;->getY()F

    move-result p1

    invoke-static {v0, v1, v2, p1}, Lorg/gioui/GioView;->onTouchExploration(JFF)V

    .line 330
    nop

    .line 335
    :goto_0
    const/4 p1, 0x1

    return p1

    nop

    :pswitch_data_0
    .packed-switch 0x7
        :pswitch_2
        :pswitch_0
        :pswitch_2
        :pswitch_1
    .end packed-switch
.end method

.method public doFrame(J)V
    .locals 3

    .line 453
    iget-wide p1, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v0, 0x0

    cmp-long v2, p1, v0

    if-eqz v2, :cond_0

    .line 454
    invoke-static {p1, p2}, Lorg/gioui/GioView;->onFrameCallback(J)V

    .line 456
    :cond_0
    return-void
.end method

.method protected fitSystemWindows(Landroid/graphics/Rect;)Z
    .locals 6

    .line 441
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-eqz v4, :cond_0

    .line 442
    iget v2, p1, Landroid/graphics/Rect;->top:I

    iget v3, p1, Landroid/graphics/Rect;->right:I

    iget v4, p1, Landroid/graphics/Rect;->bottom:I

    iget v5, p1, Landroid/graphics/Rect;->left:I

    invoke-static/range {v0 .. v5}, Lorg/gioui/GioView;->onWindowInsets(JIIII)V

    .line 444
    :cond_0
    const/4 p1, 0x1

    return p1
.end method

.method public getAccessibilityNodeProvider()Landroid/view/accessibility/AccessibilityNodeProvider;
    .locals 1

    .line 822
    new-instance v0, Lorg/gioui/GioView$2;

    invoke-direct {v0, p0}, Lorg/gioui/GioView$2;-><init>(Lorg/gioui/GioView;)V

    return-object v0
.end method

.method getDensity()I
    .locals 1

    .line 459
    invoke-virtual {p0}, Lorg/gioui/GioView;->getResources()Landroid/content/res/Resources;

    move-result-object v0

    invoke-virtual {v0}, Landroid/content/res/Resources;->getDisplayMetrics()Landroid/util/DisplayMetrics;

    move-result-object v0

    iget v0, v0, Landroid/util/DisplayMetrics;->densityDpi:I

    return v0
.end method

.method getFontScale()F
    .locals 1

    .line 463
    invoke-virtual {p0}, Lorg/gioui/GioView;->getResources()Landroid/content/res/Resources;

    move-result-object v0

    invoke-virtual {v0}, Landroid/content/res/Resources;->getConfiguration()Landroid/content/res/Configuration;

    move-result-object v0

    iget v0, v0, Landroid/content/res/Configuration;->fontScale:F

    return v0
.end method

.method hideTextInput()V
    .locals 3

    .line 437
    iget-object v0, p0, Lorg/gioui/GioView;->imm:Landroid/view/inputmethod/InputMethodManager;

    invoke-virtual {p0}, Lorg/gioui/GioView;->getWindowToken()Landroid/os/IBinder;

    move-result-object v1

    const/4 v2, 0x0

    invoke-virtual {v0, v1, v2}, Landroid/view/inputmethod/InputMethodManager;->hideSoftInputFromWindow(Landroid/os/IBinder;I)Z

    .line 438
    return-void
.end method

.method isA11yActive()Z
    .locals 1

    .line 354
    iget-object v0, p0, Lorg/gioui/GioView;->accessManager:Landroid/view/accessibility/AccessibilityManager;

    invoke-virtual {v0}, Landroid/view/accessibility/AccessibilityManager;->isEnabled()Z

    move-result v0

    return v0
.end method

.method obtainA11yEvent(II)Landroid/view/accessibility/AccessibilityEvent;
    .locals 1

    .line 347
    invoke-static {p1}, Landroid/view/accessibility/AccessibilityEvent;->obtain(I)Landroid/view/accessibility/AccessibilityEvent;

    move-result-object p1

    .line 348
    invoke-virtual {p0}, Lorg/gioui/GioView;->getContext()Landroid/content/Context;

    move-result-object v0

    invoke-virtual {v0}, Landroid/content/Context;->getPackageName()Ljava/lang/String;

    move-result-object v0

    invoke-virtual {p1, v0}, Landroid/view/accessibility/AccessibilityEvent;->setPackageName(Ljava/lang/CharSequence;)V

    .line 349
    invoke-virtual {p1, p0, p2}, Landroid/view/accessibility/AccessibilityEvent;->setSource(Landroid/view/View;I)V

    .line 350
    return-object p1
.end method

.method public onCreateInputConnection(Landroid/view/inputmethod/EditorInfo;)Landroid/view/inputmethod/InputConnection;
    .locals 4

    .line 409
    invoke-direct {p0}, Lorg/gioui/GioView;->getSnippet()Lorg/gioui/GioView$Snippet;

    move-result-object v0

    .line 410
    iget v1, p0, Lorg/gioui/GioView;->keyboardHint:I

    iput v1, p1, Landroid/view/inputmethod/EditorInfo;->inputType:I

    .line 411
    const/high16 v1, 0x12000000

    iput v1, p1, Landroid/view/inputmethod/EditorInfo;->imeOptions:I

    .line 412
    iget-wide v1, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v1, v2}, Lorg/gioui/GioView;->imeSelectionStart(J)I

    move-result v3

    invoke-static {v1, v2, v3}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v1

    iput v1, p1, Landroid/view/inputmethod/EditorInfo;->initialSelStart:I

    .line 413
    iget-wide v1, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v1, v2}, Lorg/gioui/GioView;->imeSelectionEnd(J)I

    move-result v3

    invoke-static {v1, v2, v3}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v1

    iput v1, p1, Landroid/view/inputmethod/EditorInfo;->initialSelEnd:I

    .line 414
    iget v1, p1, Landroid/view/inputmethod/EditorInfo;->initialSelStart:I

    iget v2, v0, Lorg/gioui/GioView$Snippet;->offset:I

    sub-int/2addr v1, v2

    .line 415
    iget-object v2, v0, Lorg/gioui/GioView$Snippet;->snippet:Ljava/lang/String;

    iget v3, p0, Lorg/gioui/GioView;->keyboardHint:I

    invoke-static {v2, v1, v3}, Landroid/text/TextUtils;->getCapsMode(Ljava/lang/CharSequence;II)I

    move-result v1

    iput v1, p1, Landroid/view/inputmethod/EditorInfo;->initialCapsMode:I

    .line 416
    sget v1, Landroid/os/Build$VERSION;->SDK_INT:I

    const/16 v2, 0x1e

    if-lt v1, v2, :cond_0

    .line 417
    iget-object v1, v0, Lorg/gioui/GioView$Snippet;->snippet:Ljava/lang/String;

    iget-wide v2, p0, Lorg/gioui/GioView;->nhandle:J

    iget v0, v0, Lorg/gioui/GioView$Snippet;->offset:I

    invoke-static {v2, v3, v0}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v0

    invoke-virtual {p1, v1, v0}, Landroid/view/inputmethod/EditorInfo;->setInitialSurroundingSubText(Ljava/lang/CharSequence;I)V

    .line 419
    :cond_0
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const/4 p1, -0x1

    invoke-static {v0, v1, p1, p1}, Lorg/gioui/GioView;->imeSetComposingRegion(JII)I

    .line 420
    new-instance p1, Lorg/gioui/GioView$GioInputConnection;

    const/4 v0, 0x0

    invoke-direct {p1, p0, v0}, Lorg/gioui/GioView$GioInputConnection;-><init>(Lorg/gioui/GioView;Lorg/gioui/GioView$1;)V

    return-object p1
.end method

.method public onGenericMotionEvent(Landroid/view/MotionEvent;)Z
    .locals 0

    .line 145
    invoke-direct {p0, p1}, Lorg/gioui/GioView;->dispatchMotionEvent(Landroid/view/MotionEvent;)V

    .line 146
    const/4 p1, 0x1

    return p1
.end method

.method protected onIntentEvent(Landroid/content/Intent;)V
    .locals 2

    .line 313
    if-nez p1, :cond_0

    .line 314
    return-void

    .line 316
    :cond_0
    invoke-virtual {p1}, Landroid/content/Intent;->getData()Landroid/net/Uri;

    move-result-object v0

    if-eqz v0, :cond_1

    .line 317
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-virtual {p1}, Landroid/content/Intent;->getData()Landroid/net/Uri;

    move-result-object p1

    invoke-virtual {p1}, Landroid/net/Uri;->toString()Ljava/lang/String;

    move-result-object p1

    invoke-static {v0, v1, p1}, Lorg/gioui/GioView;->onOpenURI(JLjava/lang/String;)V

    .line 319
    :cond_1
    return-void
.end method

.method public onKeyDown(ILandroid/view/KeyEvent;)Z
    .locals 7

    .line 131
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-eqz v4, :cond_0

    .line 132
    invoke-virtual {p2}, Landroid/view/KeyEvent;->getUnicodeChar()I

    move-result v3

    const/4 v4, 0x1

    invoke-virtual {p2}, Landroid/view/KeyEvent;->getEventTime()J

    move-result-wide v5

    move v2, p1

    invoke-static/range {v0 .. v6}, Lorg/gioui/GioView;->onKeyEvent(JIIZJ)V

    .line 134
    :cond_0
    const/4 p1, 0x0

    return p1
.end method

.method public onKeyUp(ILandroid/view/KeyEvent;)Z
    .locals 7

    .line 138
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-eqz v4, :cond_0

    .line 139
    invoke-virtual {p2}, Landroid/view/KeyEvent;->getUnicodeChar()I

    move-result v3

    const/4 v4, 0x0

    invoke-virtual {p2}, Landroid/view/KeyEvent;->getEventTime()J

    move-result-wide v5

    move v2, p1

    invoke-static/range {v0 .. v6}, Lorg/gioui/GioView;->onKeyEvent(JIIZJ)V

    .line 141
    :cond_0
    const/4 p1, 0x0

    return p1
.end method

.method public onTouchEvent(Landroid/view/MotionEvent;)Z
    .locals 0

    .line 152
    nop

    .line 153
    invoke-virtual {p0, p1}, Lorg/gioui/GioView;->requestUnbufferedDispatch(Landroid/view/MotionEvent;)V

    .line 156
    invoke-direct {p0, p1}, Lorg/gioui/GioView;->dispatchMotionEvent(Landroid/view/MotionEvent;)V

    .line 157
    const/4 p1, 0x1

    return p1
.end method

.method public pause()V
    .locals 5

    .line 479
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-eqz v4, :cond_0

    .line 480
    const/4 v2, 0x0

    invoke-static {v0, v1, v2}, Lorg/gioui/GioView;->onFocusChange(JZ)V

    .line 482
    :cond_0
    return-void
.end method

.method postFrameCallback()V
    .locals 1

    .line 448
    invoke-static {}, Landroid/view/Choreographer;->getInstance()Landroid/view/Choreographer;

    move-result-object v0

    invoke-virtual {v0, p0}, Landroid/view/Choreographer;->removeFrameCallback(Landroid/view/Choreographer$FrameCallback;)V

    .line 449
    invoke-static {}, Landroid/view/Choreographer;->getInstance()Landroid/view/Choreographer;

    move-result-object v0

    invoke-virtual {v0, p0}, Landroid/view/Choreographer;->postFrameCallback(Landroid/view/Choreographer$FrameCallback;)V

    .line 450
    return-void
.end method

.method restartInput()V
    .locals 1

    .line 516
    iget-object v0, p0, Lorg/gioui/GioView;->imm:Landroid/view/inputmethod/InputMethodManager;

    invoke-virtual {v0, p0}, Landroid/view/inputmethod/InputMethodManager;->restartInput(Landroid/view/View;)V

    .line 517
    return-void
.end method

.method public resume()V
    .locals 5

    .line 485
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-eqz v4, :cond_0

    .line 486
    const/4 v2, 0x1

    invoke-static {v0, v1, v2}, Lorg/gioui/GioView;->onFocusChange(JZ)V

    .line 488
    :cond_0
    return-void
.end method

.method sendA11yChange(I)V
    .locals 1

    .line 358
    iget-object v0, p0, Lorg/gioui/GioView;->accessManager:Landroid/view/accessibility/AccessibilityManager;

    invoke-virtual {v0}, Landroid/view/accessibility/AccessibilityManager;->isEnabled()Z

    move-result v0

    if-nez v0, :cond_0

    .line 359
    return-void

    .line 361
    :cond_0
    const/16 v0, 0x800

    invoke-virtual {p0, v0, p1}, Lorg/gioui/GioView;->obtainA11yEvent(II)Landroid/view/accessibility/AccessibilityEvent;

    move-result-object p1

    .line 362
    nop

    .line 363
    const/4 v0, 0x1

    invoke-virtual {p1, v0}, Landroid/view/accessibility/AccessibilityEvent;->setContentChangeTypes(I)V

    .line 365
    invoke-virtual {p0}, Lorg/gioui/GioView;->getParent()Landroid/view/ViewParent;

    move-result-object v0

    invoke-interface {v0, p0, p1}, Landroid/view/ViewParent;->requestSendAccessibilityEvent(Landroid/view/View;Landroid/view/accessibility/AccessibilityEvent;)Z

    .line 366
    return-void
.end method

.method sendA11yEvent(II)V
    .locals 1

    .line 339
    iget-object v0, p0, Lorg/gioui/GioView;->accessManager:Landroid/view/accessibility/AccessibilityManager;

    invoke-virtual {v0}, Landroid/view/accessibility/AccessibilityManager;->isEnabled()Z

    move-result v0

    if-nez v0, :cond_0

    .line 340
    return-void

    .line 342
    :cond_0
    invoke-virtual {p0, p1, p2}, Lorg/gioui/GioView;->obtainA11yEvent(II)Landroid/view/accessibility/AccessibilityEvent;

    move-result-object p1

    .line 343
    invoke-virtual {p0}, Lorg/gioui/GioView;->getParent()Landroid/view/ViewParent;

    move-result-object p2

    invoke-interface {p2, p0, p1}, Landroid/view/ViewParent;->requestSendAccessibilityEvent(Landroid/view/View;Landroid/view/accessibility/AccessibilityEvent;)Z

    .line 344
    return-void
.end method

.method setInputHint(I)V
    .locals 1

    .line 424
    iget v0, p0, Lorg/gioui/GioView;->keyboardHint:I

    if-ne p1, v0, :cond_0

    .line 425
    return-void

    .line 427
    :cond_0
    iput p1, p0, Lorg/gioui/GioView;->keyboardHint:I

    .line 428
    invoke-virtual {p0}, Lorg/gioui/GioView;->restartInput()V

    .line 429
    return-void
.end method

.method showTextInput()V
    .locals 2

    .line 432
    invoke-virtual {p0}, Lorg/gioui/GioView;->requestFocus()Z

    .line 433
    iget-object v0, p0, Lorg/gioui/GioView;->imm:Landroid/view/inputmethod/InputMethodManager;

    const/4 v1, 0x0

    invoke-virtual {v0, p0, v1}, Landroid/view/inputmethod/InputMethodManager;->showSoftInput(Landroid/view/View;I)Z

    .line 434
    return-void
.end method

.method public start()V
    .locals 5

    .line 467
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-eqz v4, :cond_0

    .line 468
    invoke-static {v0, v1}, Lorg/gioui/GioView;->onStartView(J)V

    .line 470
    :cond_0
    return-void
.end method

.method public stop()V
    .locals 5

    .line 473
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    const-wide/16 v2, 0x0

    cmp-long v4, v0, v2

    if-eqz v4, :cond_0

    .line 474
    invoke-static {v0, v1}, Lorg/gioui/GioView;->onStopView(J)V

    .line 476
    :cond_0
    return-void
.end method

.method protected unregister()V
    .locals 2

    .line 497
    const/4 v0, 0x0

    invoke-virtual {p0, v0}, Lorg/gioui/GioView;->setOnFocusChangeListener(Landroid/view/View$OnFocusChangeListener;)V

    .line 498
    invoke-virtual {p0}, Lorg/gioui/GioView;->getHolder()Landroid/view/SurfaceHolder;

    move-result-object v0

    iget-object v1, p0, Lorg/gioui/GioView;->surfCallbacks:Landroid/view/SurfaceHolder$Callback;

    invoke-interface {v0, v1}, Landroid/view/SurfaceHolder;->removeCallback(Landroid/view/SurfaceHolder$Callback;)V

    .line 499
    const-wide/16 v0, 0x0

    iput-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    .line 500
    return-void
.end method

.method updateCaret(FFFFFFFFFF)V
    .locals 8

    .line 528
    move-object v0, p0

    .line 531
    new-instance v1, Landroid/graphics/Matrix;

    invoke-direct {v1}, Landroid/graphics/Matrix;-><init>()V

    .line 532
    const/16 v2, 0x9

    new-array v2, v2, [F

    const/4 v3, 0x0

    aput p1, v2, v3

    const/4 v3, 0x1

    aput p2, v2, v3

    const/4 v3, 0x2

    aput p3, v2, v3

    const/4 v3, 0x3

    aput p4, v2, v3

    const/4 v3, 0x4

    aput p5, v2, v3

    const/4 v3, 0x5

    aput p6, v2, v3

    const/4 v3, 0x6

    const/4 v4, 0x0

    aput v4, v2, v3

    const/4 v3, 0x7

    aput v4, v2, v3

    const/16 v3, 0x8

    const/high16 v4, 0x3f800000    # 1.0f

    aput v4, v2, v3

    invoke-virtual {v1, v2}, Landroid/graphics/Matrix;->setValues([F)V

    .line 533
    invoke-virtual {p0}, Lorg/gioui/GioView;->getMatrix()Landroid/graphics/Matrix;

    move-result-object v2

    invoke-virtual {v1, v2, v1}, Landroid/graphics/Matrix;->setConcat(Landroid/graphics/Matrix;Landroid/graphics/Matrix;)Z

    .line 534
    iget-wide v2, v0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v2, v3}, Lorg/gioui/GioView;->imeSelectionStart(J)I

    move-result v2

    .line 535
    iget-wide v3, v0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v3, v4}, Lorg/gioui/GioView;->imeSelectionEnd(J)I

    move-result v3

    .line 536
    iget-wide v4, v0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v4, v5}, Lorg/gioui/GioView;->imeComposingStart(J)I

    move-result v4

    .line 537
    iget-wide v5, v0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v5, v6}, Lorg/gioui/GioView;->imeComposingEnd(J)I

    move-result v5

    .line 538
    invoke-direct {p0}, Lorg/gioui/GioView;->getSnippet()Lorg/gioui/GioView$Snippet;

    move-result-object v6

    .line 539
    nop

    .line 540
    const/4 v7, -0x1

    if-eq v4, v7, :cond_0

    .line 541
    invoke-virtual {v6, v4, v5}, Lorg/gioui/GioView$Snippet;->substringRunes(II)Ljava/lang/String;

    move-result-object v5

    goto :goto_0

    .line 540
    :cond_0
    const-string v5, ""

    .line 543
    :goto_0
    new-instance v6, Landroid/view/inputmethod/CursorAnchorInfo$Builder;

    invoke-direct {v6}, Landroid/view/inputmethod/CursorAnchorInfo$Builder;-><init>()V

    .line 544
    invoke-virtual {v6, v1}, Landroid/view/inputmethod/CursorAnchorInfo$Builder;->setMatrix(Landroid/graphics/Matrix;)Landroid/view/inputmethod/CursorAnchorInfo$Builder;

    move-result-object v1

    iget-wide v6, v0, Lorg/gioui/GioView;->nhandle:J

    .line 545
    invoke-static {v6, v7, v4}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v4

    invoke-virtual {v1, v4, v5}, Landroid/view/inputmethod/CursorAnchorInfo$Builder;->setComposingText(ILjava/lang/CharSequence;)Landroid/view/inputmethod/CursorAnchorInfo$Builder;

    move-result-object v1

    iget-wide v4, v0, Lorg/gioui/GioView;->nhandle:J

    .line 546
    invoke-static {v4, v5, v2}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v2

    iget-wide v4, v0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v4, v5, v3}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v3

    invoke-virtual {v1, v2, v3}, Landroid/view/inputmethod/CursorAnchorInfo$Builder;->setSelectionRange(II)Landroid/view/inputmethod/CursorAnchorInfo$Builder;

    move-result-object v1

    const/4 v2, 0x0

    .line 547
    move-object p1, v1

    move p2, p7

    move/from16 p3, p8

    move/from16 p4, p9

    move/from16 p5, p10

    move p6, v2

    invoke-virtual/range {p1 .. p6}, Landroid/view/inputmethod/CursorAnchorInfo$Builder;->setInsertionMarkerLocation(FFFFI)Landroid/view/inputmethod/CursorAnchorInfo$Builder;

    move-result-object v1

    .line 548
    invoke-virtual {v1}, Landroid/view/inputmethod/CursorAnchorInfo$Builder;->build()Landroid/view/inputmethod/CursorAnchorInfo;

    move-result-object v1

    .line 549
    iget-object v2, v0, Lorg/gioui/GioView;->imm:Landroid/view/inputmethod/InputMethodManager;

    invoke-virtual {v2, p0, v1}, Landroid/view/inputmethod/InputMethodManager;->updateCursorAnchorInfo(Landroid/view/View;Landroid/view/inputmethod/CursorAnchorInfo;)V

    .line 550
    return-void
.end method

.method updateSelection()V
    .locals 9

    .line 520
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v0, v1}, Lorg/gioui/GioView;->imeSelectionStart(J)I

    move-result v2

    invoke-static {v0, v1, v2}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v5

    .line 521
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v0, v1}, Lorg/gioui/GioView;->imeSelectionEnd(J)I

    move-result v2

    invoke-static {v0, v1, v2}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v6

    .line 522
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v0, v1}, Lorg/gioui/GioView;->imeComposingStart(J)I

    move-result v2

    invoke-static {v0, v1, v2}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v7

    .line 523
    iget-wide v0, p0, Lorg/gioui/GioView;->nhandle:J

    invoke-static {v0, v1}, Lorg/gioui/GioView;->imeComposingEnd(J)I

    move-result v2

    invoke-static {v0, v1, v2}, Lorg/gioui/GioView;->imeToUTF16(JI)I

    move-result v8

    .line 524
    iget-object v3, p0, Lorg/gioui/GioView;->imm:Landroid/view/inputmethod/InputMethodManager;

    move-object v4, p0

    invoke-virtual/range {v3 .. v8}, Landroid/view/inputmethod/InputMethodManager;->updateSelection(Landroid/view/View;IIII)V

    .line 525
    return-void
.end method
