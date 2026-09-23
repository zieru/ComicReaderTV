.class Lorg/gioui/GioView$Snippet;
.super Ljava/lang/Object;
.source "GioView.java"


# annotations
.annotation system Ldalvik/annotation/EnclosingClass;
    value = Lorg/gioui/GioView;
.end annotation

.annotation system Ldalvik/annotation/InnerClass;
    accessFlags = 0xa
    name = "Snippet"
.end annotation


# instance fields
.field offset:I

.field snippet:Ljava/lang/String;


# direct methods
.method private constructor <init>()V
    .locals 0

    .line 789
    invoke-direct {p0}, Ljava/lang/Object;-><init>()V

    return-void
.end method

.method synthetic constructor <init>(Lorg/gioui/GioView$1;)V
    .locals 0

    .line 789
    invoke-direct {p0}, Lorg/gioui/GioView$Snippet;-><init>()V

    return-void
.end method


# virtual methods
.method substringRunes(II)Ljava/lang/String;
    .locals 3

    .line 799
    iget v0, p0, Lorg/gioui/GioView$Snippet;->offset:I

    sub-int/2addr p1, v0

    .line 800
    sub-int/2addr p2, v0

    .line 801
    iget-object v0, p0, Lorg/gioui/GioView$Snippet;->snippet:Ljava/lang/String;

    invoke-virtual {v0}, Ljava/lang/String;->length()I

    move-result v1

    const/4 v2, 0x0

    invoke-virtual {v0, v2, v1}, Ljava/lang/String;->codePointCount(II)I

    move-result v0

    .line 802
    if-gez p1, :cond_0

    .line 803
    const/4 p1, 0x0

    .line 805
    :cond_0
    if-gez p2, :cond_1

    .line 806
    const/4 p2, 0x0

    .line 808
    :cond_1
    if-le p1, v0, :cond_2

    .line 809
    move p1, v0

    .line 811
    :cond_2
    if-le p2, v0, :cond_3

    .line 812
    goto :goto_0

    .line 811
    :cond_3
    move v0, p2

    .line 814
    :goto_0
    iget-object p2, p0, Lorg/gioui/GioView$Snippet;->snippet:Ljava/lang/String;

    .line 815
    invoke-virtual {p2, v2, p1}, Ljava/lang/String;->offsetByCodePoints(II)I

    move-result p1

    iget-object v1, p0, Lorg/gioui/GioView$Snippet;->snippet:Ljava/lang/String;

    .line 816
    invoke-virtual {v1, v2, v0}, Ljava/lang/String;->offsetByCodePoints(II)I

    move-result v0

    .line 814
    invoke-virtual {p2, p1, v0}, Ljava/lang/String;->substring(II)Ljava/lang/String;

    move-result-object p1

    return-object p1
.end method
