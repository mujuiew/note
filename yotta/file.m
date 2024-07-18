=================================================
SETTEST(input)
    S ^GlobalName(input)=""
	Q
    ;
SETTEST2(input)
    N RES
    S ^GlobalName(input,"aaaaa")=""
	Q "DONE"
    ;
=================================================
************************
****** Note !!!! *******
************************
zed "filename"                  เขียนไฟล์
zl "filename"                   compiler ไฟล์.o                                
cd /ydbdir/rtns                 path วางไฟล์ TEST.m
cd obj/                         path วางไฟล์ TEST.o

d funcName^fileName("xx")       กรณีไม่มี return ค่า           : d SETTEST^TEST("11111")
w $$funcName^fileName("xx")     กรณีมี return ค่า             : w $$SETTEST^TEST("11111")
=================================================
[ydbadm@host rtns]$ cd /ydbdir
[ydbadm@host ydbdir]$ . ./ydbenv
[ydbadm@host ydbdir]$ ydb

YDB> zed "TEST"
YDB> zl "TEST"

YDB> d SETTEST^TEST("11111")
YDB> d SETTEST^TEST("626005020001")

YDB> zwr ^TEST
^TEST(11111)=""
^TEST(626005020001)=""

YDB> w $$SETTEST^TEST("11111")
DONE


YDB> w $$SETTEST3^TESTFUNC("2019-08-30")
