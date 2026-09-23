.class Lorg/gioui/GioView$GioInputConnection;
.super Ljava/lang/Object;
.source "GioView.java"

# interfaces
.implements Landroid/view/inputmethod/InputConnection;


# annotations
.annotation system Ldalvik/annotation/EnclosingClass;
    value = Lorg/gioui/GioView;
.end annotation

.annotation system Ldalvik/annotation/InnerClass;
    accessFlags = 0x2
    name = "GioInputConnection"
.end annotation


# instance fields
.field private batchDepth:I

.field final synthetic this$0:Lorg/gioui/GioView;


# direct methods
.method private constructor <init>(Lorg/gioui/GioView;)V
    .locals 0
    .annotation system Ldalvik/annotation/MethodParameters;
        accessFlags = {
            0x1010
        }
        names = {
            null
        }
    .end annotation

    .line 587
    iput-object p1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-direct {p0}, Ljava/lang/Object;-><init>()V

    return-void
.end method

.method synthetic constructor <init>(Lorg/gioui/GioView;Lorg/gioui/GioView$1;)V
    .locals 0

    .line 587
    invoke-direct {p0, p1}, Lorg/gioui/GioView$GioInputConnection;-><init>(Lorg/gioui/GioView;)V

    return-void
.end method


# virtual methods
.method public beginBatchEdit()Z
    .locals 2

    .line 591
    iget v0, p0, Lorg/gioui/GioView$GioInputConnection;->batchDepth:I

    const/4 v1, 0x1

    add-int/2addr v0, v1

    iput v0, p0, Lorg/gioui/GioView$GioInputConnection;->batchDepth:I

    .line 592
    return v1
.end method

.method public clearMetaKeyStates(I)Z
    .locals 0

    .line 601
    const/4 p1, 0x0

    return p1
.end method

.method public closeConnection()V
    .locals 0

    .line 749
    return-void
.end method

.method public commitCompletion(Landroid/view/inputmethod/CompletionInfo;)Z
    .locals 0

    .line 605
    const/4 p1, 0x0

    return p1
.end method

.method public commitContent(Landroid/view/inputmethod/InputContentInfo;ILandroid/os/Bundle;)Z
    .locals 0

    .line 756
    const/4 p1, 0x0

    return p1
.end method

.method public commitCorrection(Landroid/view/inputmethod/CorrectionInfo;)Z
    .locals 0

    .line 609
    const/4 p1, 0x0

    return p1
.end method

.method public commitText(Ljava/lang/CharSequence;I)Z
    .locals 0

    .line 613
    invoke-virtual {p0, p1, p2}, Lorg/gioui/GioView$GioInputConnection;->setComposingText(Ljava/lang/CharSequence;I)Z

    .line 614
    invoke-virtual {p0}, Lorg/gioui/GioView$GioInputConnection;->finishComposingText()Z

    move-result p1

    return p1
.end method

.method public deleteSurroundingText(II)Z
    .locals 6

    .line 619
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1}, Lorg/gioui/GioView;->access$400(J)I

    move-result v0

    .line 620
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$500(J)I

    move-result v1

    .line 621
    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    iget-object v4, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v4}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v4

    invoke-static {v4, v5, v0}, Lorg/gioui/GioView;->access$600(JI)I

    move-result v4

    sub-int/2addr v4, p1

    invoke-static {v2, v3, v4}, Lorg/gioui/GioView;->access$700(JI)I

    move-result p1

    sub-int/2addr v0, p1

    .line 622
    iget-object p1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    iget-object p1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v4

    invoke-static {v4, v5, v1}, Lorg/gioui/GioView;->access$600(JI)I

    move-result p1

    sub-int/2addr p1, p2

    invoke-static {v2, v3, p1}, Lorg/gioui/GioView;->access$700(JI)I

    move-result p1

    sub-int/2addr v1, p1

    .line 623
    invoke-virtual {p0, v0, v1}, Lorg/gioui/GioView$GioInputConnection;->deleteSurroundingTextInCodePoints(II)Z

    move-result p1

    return p1
.end method

.method public deleteSurroundingTextInCodePoints(II)Z
    .locals 4

    .line 760
    const-string v0, ""

    if-lez p2, :cond_0

    .line 761
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$500(J)I

    move-result v1

    .line 762
    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    add-int/2addr p2, v1

    invoke-static {v2, v3, v1, p2, v0}, Lorg/gioui/GioView;->access$1400(JIILjava/lang/String;)I

    .line 764
    :cond_0
    if-lez p1, :cond_1

    .line 765
    iget-object p2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$400(J)I

    move-result p2

    .line 766
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    sub-int p1, p2, p1

    invoke-static {v1, v2, p1, p2, v0}, Lorg/gioui/GioView;->access$1400(JIILjava/lang/String;)I

    .line 768
    :cond_1
    const/4 p1, 0x1

    return p1
.end method

.method public endBatchEdit()Z
    .locals 2

    .line 596
    iget v0, p0, Lorg/gioui/GioView$GioInputConnection;->batchDepth:I

    const/4 v1, 0x1

    sub-int/2addr v0, v1

    iput v0, p0, Lorg/gioui/GioView$GioInputConnection;->batchDepth:I

    .line 597
    if-lez v0, :cond_0

    goto :goto_0

    :cond_0
    const/4 v1, 0x0

    :goto_0
    return v1
.end method

.method public finishComposingText()Z
    .locals 3

    .line 627
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    const/4 v2, -0x1

    invoke-static {v0, v1, v2, v2}, Lorg/gioui/GioView;->access$800(JII)I

    .line 628
    const/4 v0, 0x1

    return v0
.end method

.method public getCursorCapsMode(I)I
    .locals 5

    .line 632
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$900(Lorg/gioui/GioView;)Lorg/gioui/GioView$Snippet;

    move-result-object v0

    .line 633
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$400(J)I

    move-result v1

    .line 634
    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    iget v4, v0, Lorg/gioui/GioView$Snippet;->offset:I

    sub-int/2addr v1, v4

    invoke-static {v2, v3, v1}, Lorg/gioui/GioView;->access$600(JI)I

    move-result v1

    .line 635
    if-ltz v1, :cond_1

    iget-object v2, v0, Lorg/gioui/GioView$Snippet;->snippet:Ljava/lang/String;

    invoke-virtual {v2}, Ljava/lang/String;->length()I

    move-result v2

    if-le v1, v2, :cond_0

    goto :goto_0

    .line 638
    :cond_0
    iget-object v0, v0, Lorg/gioui/GioView$Snippet;->snippet:Ljava/lang/String;

    invoke-static {v0, v1, p1}, Landroid/text/TextUtils;->getCapsMode(Ljava/lang/CharSequence;II)I

    move-result p1

    return p1

    .line 636
    :cond_1
    :goto_0
    const/4 p1, 0x0

    return p1
.end method

.method public getExtractedText(Landroid/view/inputmethod/ExtractedTextRequest;I)Landroid/view/inputmethod/ExtractedText;
    .locals 0

    .line 642
    const/4 p1, 0x0

    return-object p1
.end method

.method public getHandler()Landroid/os/Handler;
    .locals 1

    .line 752
    const/4 v0, 0x0

    return-object v0
.end method

.method public getSelectedText(I)Ljava/lang/CharSequence;
    .locals 3

    .line 646
    iget-object p1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$900(Lorg/gioui/GioView;)Lorg/gioui/GioView$Snippet;

    move-result-object p1

    .line 647
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1}, Lorg/gioui/GioView;->access$400(J)I

    move-result v0

    .line 648
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$500(J)I

    move-result v1

    .line 649
    invoke-virtual {p1, v0, v1}, Lorg/gioui/GioView$Snippet;->substringRunes(II)Ljava/lang/String;

    move-result-object p1

    .line 650
    return-object p1
.end method

.method public getSurroundingText(III)Landroid/view/inputmethod/SurroundingText;
    .locals 4

    .line 772
    iget-object p3, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p3}, Lorg/gioui/GioView;->access$900(Lorg/gioui/GioView;)Lorg/gioui/GioView$Snippet;

    move-result-object p3

    .line 773
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1}, Lorg/gioui/GioView;->access$400(J)I

    move-result v0

    .line 774
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$500(J)I

    move-result v1

    .line 776
    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    sub-int p1, v0, p1

    add-int/2addr p2, v1

    invoke-static {v2, v3, p1, p2}, Lorg/gioui/GioView;->access$1000(JII)V

    .line 777
    new-instance p1, Landroid/view/inputmethod/SurroundingText;

    iget-object p2, p3, Lorg/gioui/GioView$Snippet;->snippet:Ljava/lang/String;

    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    invoke-static {v2, v3, v0}, Lorg/gioui/GioView;->access$600(JI)I

    move-result v0

    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    invoke-static {v2, v3, v1}, Lorg/gioui/GioView;->access$600(JI)I

    move-result v1

    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    iget p3, p3, Lorg/gioui/GioView$Snippet;->offset:I

    invoke-static {v2, v3, p3}, Lorg/gioui/GioView;->access$600(JI)I

    move-result p3

    invoke-direct {p1, p2, v0, v1, p3}, Landroid/view/inputmethod/SurroundingText;-><init>(Ljava/lang/CharSequence;III)V

    return-object p1
.end method

.method public getTextAfterCursor(II)Ljava/lang/CharSequence;
    .locals 6

    .line 654
    iget-object p2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p2}, Lorg/gioui/GioView;->access$900(Lorg/gioui/GioView;)Lorg/gioui/GioView$Snippet;

    move-result-object p2

    .line 655
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1}, Lorg/gioui/GioView;->access$400(J)I

    move-result v0

    .line 656
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$500(J)I

    move-result v1

    .line 659
    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    sub-int/2addr v0, p1

    add-int v4, v1, p1

    invoke-static {v2, v3, v0, v4}, Lorg/gioui/GioView;->access$1000(JII)V

    .line 660
    nop

    .line 661
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v4

    invoke-static {v4, v5, v1}, Lorg/gioui/GioView;->access$600(JI)I

    move-result v0

    add-int/2addr v0, p1

    invoke-static {v2, v3, v0}, Lorg/gioui/GioView;->access$700(JI)I

    move-result p1

    .line 662
    invoke-virtual {p2, v1, p1}, Lorg/gioui/GioView$Snippet;->substringRunes(II)Ljava/lang/String;

    move-result-object p1

    .line 663
    return-object p1
.end method

.method public getTextBeforeCursor(II)Ljava/lang/CharSequence;
    .locals 5

    .line 667
    iget-object p2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p2}, Lorg/gioui/GioView;->access$900(Lorg/gioui/GioView;)Lorg/gioui/GioView$Snippet;

    move-result-object p2

    .line 668
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1}, Lorg/gioui/GioView;->access$400(J)I

    move-result v0

    .line 669
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$500(J)I

    move-result v1

    .line 672
    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    sub-int v4, v0, p1

    add-int/2addr v1, p1

    invoke-static {v2, v3, v4, v1}, Lorg/gioui/GioView;->access$1000(JII)V

    .line 673
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    iget-object v3, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v3}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v3

    invoke-static {v3, v4, v0}, Lorg/gioui/GioView;->access$600(JI)I

    move-result v3

    sub-int/2addr v3, p1

    invoke-static {v1, v2, v3}, Lorg/gioui/GioView;->access$700(JI)I

    move-result p1

    .line 674
    nop

    .line 675
    invoke-virtual {p2, p1, v0}, Lorg/gioui/GioView$Snippet;->substringRunes(II)Ljava/lang/String;

    move-result-object p1

    .line 676
    return-object p1
.end method

.method public performContextMenuAction(I)Z
    .locals 0

    .line 680
    const/4 p1, 0x0

    return p1
.end method

.method public performEditorAction(I)Z
    .locals 9

    .line 684
    invoke-static {}, Landroid/os/SystemClock;->uptimeMillis()J

    move-result-wide v7

    .line 686
    iget-object p1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    const/16 v2, 0x42

    const/16 v3, 0xa

    const/4 v4, 0x1

    move-wide v5, v7

    invoke-static/range {v0 .. v6}, Lorg/gioui/GioView;->access$1100(JIIZJ)V

    .line 687
    iget-object p1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    const/4 v4, 0x0

    invoke-static/range {v0 .. v6}, Lorg/gioui/GioView;->access$1100(JIIZJ)V

    .line 688
    const/4 p1, 0x1

    return p1
.end method

.method public performPrivateCommand(Ljava/lang/String;Landroid/os/Bundle;)Z
    .locals 0

    .line 692
    const/4 p1, 0x0

    return p1
.end method

.method public reportFullscreenMode(Z)Z
    .locals 0

    .line 696
    const/4 p1, 0x0

    return p1
.end method

.method public requestCursorUpdates(I)Z
    .locals 0

    .line 745
    const/4 p1, 0x1

    return p1
.end method

.method public sendKeyEvent(Landroid/view/KeyEvent;)Z
    .locals 9

    .line 700
    invoke-virtual {p1}, Landroid/view/KeyEvent;->getAction()I

    move-result v0

    const/4 v1, 0x1

    if-nez v0, :cond_0

    const/4 v6, 0x1

    goto :goto_0

    :cond_0
    const/4 v0, 0x0

    const/4 v6, 0x0

    .line 701
    :goto_0
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    invoke-virtual {p1}, Landroid/view/KeyEvent;->getKeyCode()I

    move-result v4

    invoke-virtual {p1}, Landroid/view/KeyEvent;->getUnicodeChar()I

    move-result v5

    invoke-virtual {p1}, Landroid/view/KeyEvent;->getEventTime()J

    move-result-wide v7

    invoke-static/range {v2 .. v8}, Lorg/gioui/GioView;->access$1100(JIIZJ)V

    .line 702
    return v1
.end method

.method public setComposingRegion(II)Z
    .locals 2

    .line 706
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1, p1}, Lorg/gioui/GioView;->access$700(JI)I

    move-result p1

    .line 707
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1, p2}, Lorg/gioui/GioView;->access$700(JI)I

    move-result p2

    .line 708
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1, p1, p2}, Lorg/gioui/GioView;->access$800(JII)I

    .line 709
    const/4 p1, 0x1

    return p1
.end method

.method public setComposingText(Ljava/lang/CharSequence;I)Z
    .locals 6

    .line 713
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1}, Lorg/gioui/GioView;->access$1200(J)I

    move-result v0

    .line 714
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$1300(J)I

    move-result v1

    .line 715
    const/4 v2, -0x1

    if-eq v0, v2, :cond_0

    if-ne v1, v2, :cond_1

    .line 716
    :cond_0
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1}, Lorg/gioui/GioView;->access$400(J)I

    move-result v0

    .line 717
    iget-object v1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v1

    invoke-static {v1, v2}, Lorg/gioui/GioView;->access$500(J)I

    move-result v1

    .line 719
    :cond_1
    invoke-interface {p1}, Ljava/lang/CharSequence;->toString()Ljava/lang/String;

    move-result-object p1

    .line 720
    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    invoke-static {v2, v3, v0, v1, p1}, Lorg/gioui/GioView;->access$1400(JIILjava/lang/String;)I

    .line 721
    nop

    .line 722
    const/4 v1, 0x0

    invoke-virtual {p1}, Ljava/lang/String;->length()I

    move-result v2

    invoke-virtual {p1, v1, v2}, Ljava/lang/String;->codePointCount(II)I

    move-result p1

    .line 723
    if-lez p2, :cond_2

    .line 724
    add-int v1, v0, p1

    .line 725
    add-int/lit8 p2, p2, -0x1

    goto :goto_0

    .line 723
    :cond_2
    move v1, v0

    .line 727
    :goto_0
    iget-object v2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    add-int/2addr p1, v0

    invoke-static {v2, v3, v0, p1}, Lorg/gioui/GioView;->access$800(JII)I

    .line 730
    iget-object p1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$900(Lorg/gioui/GioView;)Lorg/gioui/GioView$Snippet;

    .line 731
    iget-object p1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v2

    iget-object p1, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p1}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v4

    invoke-static {v4, v5, v1}, Lorg/gioui/GioView;->access$600(JI)I

    move-result p1

    add-int/2addr p1, p2

    invoke-static {v2, v3, p1}, Lorg/gioui/GioView;->access$700(JI)I

    move-result p1

    .line 732
    iget-object p2, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {p2}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1, p1, p1}, Lorg/gioui/GioView;->access$1500(JII)I

    .line 733
    const/4 p1, 0x1

    return p1
.end method

.method public setSelection(II)Z
    .locals 2

    .line 737
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1, p1}, Lorg/gioui/GioView;->access$700(JI)I

    move-result p1

    .line 738
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1, p2}, Lorg/gioui/GioView;->access$700(JI)I

    move-result p2

    .line 739
    iget-object v0, p0, Lorg/gioui/GioView$GioInputConnection;->this$0:Lorg/gioui/GioView;

    invoke-static {v0}, Lorg/gioui/GioView;->access$000(Lorg/gioui/GioView;)J

    move-result-wide v0

    invoke-static {v0, v1, p1, p2}, Lorg/gioui/GioView;->access$1500(JII)I

    .line 740
    const/4 p1, 0x1

    return p1
.end method
