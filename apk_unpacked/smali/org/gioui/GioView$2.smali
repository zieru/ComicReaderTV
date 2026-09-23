.class Lorg/gioui/GioView$2;
.super Landroid/view/accessibility/AccessibilityNodeProvider;
.source "GioView.java"


# annotations
.annotation system Ldalvik/annotation/EnclosingMethod;
    value = Lorg/gioui/GioView;->getAccessibilityNodeProvider()Landroid/view/accessibility/AccessibilityNodeProvider;
.end annotation

.annotation system Ldalvik/annotation/InnerClass;
    accessFlags = 0x0
    name = null
.end annotation


# instance fields
.field private final screenOff:[I

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

    .line 822
    iput-object p1, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    invoke-direct {p0}, Landroid/view/accessibility/AccessibilityNodeProvider;-><init>()V

    .line 823
    const/4 p1, 0x2

    new-array p1, p1, [I

    iput-object p1, p0, Lorg/gioui/GioView$2;->screenOff:[I

    return-void
.end method


# virtual methods
.method public createAccessibilityNodeInfo(I)Landroid/view/accessibility/AccessibilityNodeInfo;
    .locals 9

    .line 826
    nop

    .line 827
    const/4 v0, -0x1

    const/4 v1, 0x1

    if-ne p1, v0, :cond_0

    .line 828
    iget-object v0, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Landroid/view/accessibility/AccessibilityNodeInfo;->obtain(Landroid/view/View;)Landroid/view/accessibility/AccessibilityNodeInfo;

    move-result-object v0

    .line 829
    iget-object v2, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    invoke-virtual {v2, v0}, Lorg/gioui/GioView;->onInitializeAccessibilityNodeInfo(Landroid/view/accessibility/AccessibilityNodeInfo;)V

    move-object v8, v0

    goto :goto_0

    .line 831
    :cond_0
    iget-object v0, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    invoke-static {v0, p1}, Landroid/view/accessibility/AccessibilityNodeInfo;->obtain(Landroid/view/View;I)Landroid/view/accessibility/AccessibilityNodeInfo;

    move-result-object v0

    .line 832
    iget-object v2, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    invoke-virtual {v2}, Lorg/gioui/GioView;->getContext()Landroid/content/Context;

    move-result-object v2

    invoke-virtual {v2}, Landroid/content/Context;->getPackageName()Ljava/lang/String;

    move-result-object v2

    invoke-virtual {v0, v2}, Landroid/view/accessibility/AccessibilityNodeInfo;->setPackageName(Ljava/lang/CharSequence;)V

    .line 833
    invoke-virtual {v0, v1}, Landroid/view/accessibility/AccessibilityNodeInfo;->setVisibleToUser(Z)V

    move-object v8, v0

    .line 835
    :goto_0
    iget-object v0, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    iget-object v2, p0, Lorg/gioui/GioView$2;->screenOff:[I

    invoke-virtual {v0, v2}, Lorg/gioui/GioView;->getLocationOnScreen([I)V

    .line 836
    iget-object v0, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v3

    iget-object v0, p0, Lorg/gioui/GioView$2;->screenOff:[I

    const/4 v2, 0x0

    aget v6, v0, v2

    aget v7, v0, v1

    move v5, p1

    invoke-static/range {v3 .. v8}, Lorg/gioui/GioView;->access$1700(JIIILandroid/view/accessibility/AccessibilityNodeInfo;)Landroid/view/accessibility/AccessibilityNodeInfo;

    move-result-object p1

    .line 837
    return-object p1
.end method

.method public performAction(IILandroid/os/Bundle;)Z
    .locals 2

    .line 841
    const/4 v0, -0x1

    if-ne p1, v0, :cond_0

    .line 842
    iget-object p1, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    invoke-virtual {p1, p2, p3}, Lorg/gioui/GioView;->performAccessibilityAction(ILandroid/os/Bundle;)Z

    move-result p1

    return p1

    .line 844
    :cond_0
    const/4 p3, 0x1

    sparse-switch p2, :sswitch_data_0

    .line 854
    const/4 p1, 0x0

    return p1

    .line 850
    :sswitch_0
    iget-object p2, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    invoke-static {p2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1, p1}, Lorg/gioui/GioView;->access$1900(JI)V

    .line 851
    iget-object p2, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    const/high16 v0, 0x10000

    invoke-virtual {p2, v0, p1}, Lorg/gioui/GioView;->sendA11yEvent(II)V

    .line 852
    return p3

    .line 846
    :sswitch_1
    iget-object p2, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    invoke-static {p2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1, p1}, Lorg/gioui/GioView;->access$1800(JI)V

    .line 847
    iget-object p2, p0, Lorg/gioui/GioView$2;->this$0:Lorg/gioui/GioView;

    const v0, 0x8000

    invoke-virtual {p2, v0, p1}, Lorg/gioui/GioView;->sendA11yEvent(II)V

    .line 848
    return p3

    nop

    :sswitch_data_0
    .sparse-switch
        0x40 -> :sswitch_1
        0x80 -> :sswitch_0
    .end sparse-switch
.end method
