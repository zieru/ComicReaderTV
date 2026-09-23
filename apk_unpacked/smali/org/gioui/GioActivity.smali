.class public final Lorg/gioui/GioActivity;
.super Landroid/app/Activity;
.source "GioActivity.java"


# instance fields
.field public layer:Landroid/widget/FrameLayout;

.field private view:Lorg/gioui/GioView;


# direct methods
.method public constructor <init>()V
    .locals 0

    .line 14
    invoke-direct {p0}, Landroid/app/Activity;-><init>()V

    return-void
.end method


# virtual methods
.method public onBackPressed()V
    .locals 1

    .line 72
    iget-object v0, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {v0}, Lorg/gioui/GioView;->backPressed()Z

    move-result v0

    if-nez v0, :cond_0

    .line 73
    invoke-super {p0}, Landroid/app/Activity;->onBackPressed()V

    .line 74
    :cond_0
    return-void
.end method

.method public onConfigurationChanged(Landroid/content/res/Configuration;)V
    .locals 0

    .line 62
    invoke-super {p0, p1}, Landroid/app/Activity;->onConfigurationChanged(Landroid/content/res/Configuration;)V

    .line 63
    iget-object p1, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {p1}, Lorg/gioui/GioView;->configurationChanged()V

    .line 64
    return-void
.end method

.method public onCreate(Landroid/os/Bundle;)V
    .locals 2

    .line 19
    invoke-super {p0, p1}, Landroid/app/Activity;->onCreate(Landroid/os/Bundle;)V

    .line 21
    new-instance p1, Landroid/widget/FrameLayout;

    invoke-direct {p1, p0}, Landroid/widget/FrameLayout;-><init>(Landroid/content/Context;)V

    iput-object p1, p0, Lorg/gioui/GioActivity;->layer:Landroid/widget/FrameLayout;

    .line 22
    new-instance p1, Lorg/gioui/GioView;

    invoke-direct {p1, p0}, Lorg/gioui/GioView;-><init>(Landroid/content/Context;)V

    iput-object p1, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    .line 24
    new-instance v0, Landroid/widget/FrameLayout$LayoutParams;

    const/4 v1, -0x1

    invoke-direct {v0, v1, v1}, Landroid/widget/FrameLayout$LayoutParams;-><init>(II)V

    invoke-virtual {p1, v0}, Lorg/gioui/GioView;->setLayoutParams(Landroid/view/ViewGroup$LayoutParams;)V

    .line 28
    iget-object p1, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    const/4 v0, 0x1

    invoke-virtual {p1, v0}, Lorg/gioui/GioView;->setFocusable(Z)V

    .line 29
    iget-object p1, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {p1, v0}, Lorg/gioui/GioView;->setFocusableInTouchMode(Z)V

    .line 31
    iget-object p1, p0, Lorg/gioui/GioActivity;->layer:Landroid/widget/FrameLayout;

    iget-object v0, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {p1, v0}, Landroid/widget/FrameLayout;->addView(Landroid/view/View;)V

    .line 32
    iget-object p1, p0, Lorg/gioui/GioActivity;->layer:Landroid/widget/FrameLayout;

    invoke-virtual {p0, p1}, Lorg/gioui/GioActivity;->setContentView(Landroid/view/View;)V

    .line 33
    invoke-virtual {p0}, Lorg/gioui/GioActivity;->getIntent()Landroid/content/Intent;

    move-result-object p1

    invoke-virtual {p0, p1}, Lorg/gioui/GioActivity;->onNewIntent(Landroid/content/Intent;)V

    .line 34
    return-void
.end method

.method public onDestroy()V
    .locals 1

    .line 37
    iget-object v0, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {v0}, Lorg/gioui/GioView;->destroy()V

    .line 38
    invoke-super {p0}, Landroid/app/Activity;->onDestroy()V

    .line 39
    return-void
.end method

.method public onLowMemory()V
    .locals 0

    .line 67
    invoke-super {p0}, Landroid/app/Activity;->onLowMemory()V

    .line 68
    invoke-static {}, Lorg/gioui/GioView;->onLowMemory()V

    .line 69
    return-void
.end method

.method protected onNewIntent(Landroid/content/Intent;)V
    .locals 1

    .line 77
    invoke-super {p0, p1}, Landroid/app/Activity;->onNewIntent(Landroid/content/Intent;)V

    .line 78
    iget-object v0, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {v0, p1}, Lorg/gioui/GioView;->onIntentEvent(Landroid/content/Intent;)V

    .line 79
    return-void
.end method

.method public onPause()V
    .locals 1

    .line 52
    invoke-super {p0}, Landroid/app/Activity;->onPause()V

    .line 53
    iget-object v0, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {v0}, Lorg/gioui/GioView;->pause()V

    .line 54
    return-void
.end method

.method public onResume()V
    .locals 1

    .line 57
    invoke-super {p0}, Landroid/app/Activity;->onResume()V

    .line 58
    iget-object v0, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {v0}, Lorg/gioui/GioView;->resume()V

    .line 59
    return-void
.end method

.method public onStart()V
    .locals 1

    .line 42
    invoke-super {p0}, Landroid/app/Activity;->onStart()V

    .line 43
    iget-object v0, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {v0}, Lorg/gioui/GioView;->start()V

    .line 44
    return-void
.end method

.method public onStop()V
    .locals 1

    .line 47
    iget-object v0, p0, Lorg/gioui/GioActivity;->view:Lorg/gioui/GioView;

    invoke-virtual {v0}, Lorg/gioui/GioView;->stop()V

    .line 48
    invoke-super {p0}, Landroid/app/Activity;->onStop()V

    .line 49
    return-void
.end method
