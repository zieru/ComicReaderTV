.class Lorg/gioui/GioView$1;
.super Ljava/lang/Object;
.source "GioView.java"

# interfaces
.implements Landroid/view/SurfaceHolder$Callback;


# annotations
.annotation system Ldalvik/annotation/EnclosingMethod;
    value = Lorg/gioui/GioView;-><init>(Landroid/content/Context;Landroid/util/AttributeSet;)V
.end annotation

.annotation system Ldalvik/annotation/InnerClass;
    accessFlags = 0x0
    name = null
.end annotation


# instance fields
.field final synthetic this$0:Lorg/gioui/GioView;


# direct methods
.method constructor <init>(Lorg/gioui/GioView;)V
    .locals 0
    .annotation system Ldalvik/annotation/MethodParameters;
        accessFlags = {
            0x8010
        }
        names = {
            null
        }
    .end annotation

    .line 116
    iput-object p1, p0, Lorg/gioui/GioView$1;->this$0:Lorg/gioui/GioView;

    invoke-direct {p0}, Ljava/lang/Object;-><init>()V

    return-void
.end method


# virtual methods
.method public surfaceChanged(Landroid/view/SurfaceHolder;III)V
    .locals 0

    .line 121
    iget-object p1, p0, Lorg/gioui/GioView$1;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide p1

    iget-object p3, p0, Lorg/gioui/GioView$1;->this$0:Lorg/gioui/GioView;

    invoke-virtual {p3}, Lorg/gioui/GioView;->getHolder()Landroid/view/SurfaceHolder;

    move-result-object p3

    invoke-interface {p3}, Landroid/view/SurfaceHolder;->getSurface()Landroid/view/Surface;

    move-result-object p3

    invoke-static {p1, p2, p3}, Lorg/gioui/GioView;->access$100(JLandroid/view/Surface;)V

    .line 122
    return-void
.end method

.method public surfaceCreated(Landroid/view/SurfaceHolder;)V
    .locals 0

    .line 119
    return-void
.end method

.method public surfaceDestroyed(Landroid/view/SurfaceHolder;)V
    .locals 2

    .line 124
    iget-object p1, p0, Lorg/gioui/GioView$1;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1}, Lorg/gioui/GioView;->access$200(J)V

    .line 125
    return-void
.end method
