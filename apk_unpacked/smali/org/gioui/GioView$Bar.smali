.class final enum Lorg/gioui/GioView$Bar;
.super Ljava/lang/Enum;
.source "GioView.java"


# annotations
.annotation system Ldalvik/annotation/EnclosingClass;
    value = Lorg/gioui/GioView;
.end annotation

.annotation system Ldalvik/annotation/InnerClass;
    accessFlags = 0x401a
    name = "Bar"
.end annotation

.annotation system Ldalvik/annotation/Signature;
    value = {
        "Ljava/lang/Enum<",
        "Lorg/gioui/GioView$Bar;",
        ">;"
    }
.end annotation


# static fields
.field private static final synthetic $VALUES:[Lorg/gioui/GioView$Bar;

.field public static final enum NAVIGATION:Lorg/gioui/GioView$Bar;

.field public static final enum STATUS:Lorg/gioui/GioView$Bar;


# direct methods
.method private static synthetic $values()[Lorg/gioui/GioView$Bar;
    .locals 3

    .line 191
    const/4 v0, 0x2

    new-array v0, v0, [Lorg/gioui/GioView$Bar;

    const/4 v1, 0x0

    sget-object v2, Lorg/gioui/GioView$Bar;->NAVIGATION:Lorg/gioui/GioView$Bar;

    aput-object v2, v0, v1

    const/4 v1, 0x1

    sget-object v2, Lorg/gioui/GioView$Bar;->STATUS:Lorg/gioui/GioView$Bar;

    aput-object v2, v0, v1

    return-object v0
.end method

.method static constructor <clinit>()V
    .locals 3

    .line 192
    new-instance v0, Lorg/gioui/GioView$Bar;

    const-string v1, "NAVIGATION"

    const/4 v2, 0x0

    invoke-direct {v0, v1, v2}, Lorg/gioui/GioView$Bar;-><init>(Ljava/lang/String;I)V

    sput-object v0, Lorg/gioui/GioView$Bar;->NAVIGATION:Lorg/gioui/GioView$Bar;

    .line 193
    new-instance v0, Lorg/gioui/GioView$Bar;

    const-string v1, "STATUS"

    const/4 v2, 0x1

    invoke-direct {v0, v1, v2}, Lorg/gioui/GioView$Bar;-><init>(Ljava/lang/String;I)V

    sput-object v0, Lorg/gioui/GioView$Bar;->STATUS:Lorg/gioui/GioView$Bar;

    .line 191
    invoke-static {}, Lorg/gioui/GioView$Bar;->$values()[Lorg/gioui/GioView$Bar;

    move-result-object v0

    sput-object v0, Lorg/gioui/GioView$Bar;->$VALUES:[Lorg/gioui/GioView$Bar;

    return-void
.end method

.method private constructor <init>(Ljava/lang/String;I)V
    .locals 0
    .annotation system Ldalvik/annotation/MethodParameters;
        accessFlags = {
            0x1000,
            0x1000
        }
        names = {
            null,
            null
        }
    .end annotation

    .annotation system Ldalvik/annotation/Signature;
        value = {
            "()V"
        }
    .end annotation

    .line 191
    invoke-direct {p0, p1, p2}, Ljava/lang/Enum;-><init>(Ljava/lang/String;I)V

    return-void
.end method

.method public static valueOf(Ljava/lang/String;)Lorg/gioui/GioView$Bar;
    .locals 1
    .annotation system Ldalvik/annotation/MethodParameters;
        accessFlags = {
            0x8000
        }
        names = {
            null
        }
    .end annotation

    .line 191
    const-class v0, Lorg/gioui/GioView$Bar;

    invoke-static {v0, p0}, Ljava/lang/Enum;->valueOf(Ljava/lang/Class;Ljava/lang/String;)Ljava/lang/Enum;

    move-result-object p0

    check-cast p0, Lorg/gioui/GioView$Bar;

    return-object p0
.end method

.method public static values()[Lorg/gioui/GioView$Bar;
    .locals 1

    .line 191
    sget-object v0, Lorg/gioui/GioView$Bar;->$VALUES:[Lorg/gioui/GioView$Bar;

    invoke-virtual {v0}, [Lorg/gioui/GioView$Bar;->clone()Ljava/lang/Object;

    move-result-object v0

    check-cast v0, [Lorg/gioui/GioView$Bar;

    return-object v0
.end method
