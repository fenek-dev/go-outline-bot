package markup

func Init() {

	Menu.Reply(
		Menu.Row(
			InfoOpenBtn,
		),
	)

	Info.Inline(
		Info.Row(ClientListOpenBtn),
		Info.Row(InfoClose),
	)

	ClientList.Inline(
		ClientList.Row(ClientListIOS),
		ClientList.Row(ClientListAndroid),
		ClientList.Row(ClientListWindows),
		ClientList.Row(ClientListMacOS),
		ClientList.Row(ClientListBack),
	)

	IOSList.Inline(
		IOSList.Row(IOSList.URL("Outline", "https://apps.apple.com/ru/app/outline-app/id1356177741")),
		IOSList.Row(IOSList.URL("Hiddify-next (github)", "https://github.com/hiddify/hiddify-next/releases")),
		IOSList.Row(IOSList.URL("Shadowrocket", "https://apps.apple.com/us/app/shadowrocket/id932747118")),
		IOSList.Row(IOSList.URL("FoXray", "https://apps.apple.com/us/app/foxray/id6448898396")),
		IOSList.Row(IOSListBackBtn),
	)

	AndroidList.Inline(
		AndroidList.Row(AndroidList.URL("Outline", "https://play.google.com/store/apps/details?id=org.outline.android.client")),
		AndroidList.Row(AndroidList.URL("Hiddify-next", "https://play.google.com/store/apps/details?id=app.hiddify.com")),
		AndroidList.Row(AndroidList.URL("Hiddify-next (github)", "https://github.com/hiddify/hiddify-next/releases")),
		AndroidList.Row(AndroidList.URL("v2rayNG", "https://play.google.com/store/apps/details?id=com.v2ray.ang")),
		AndroidList.Row(AndroidList.URL("NekoBox", "https://github.com/MatsuriDayo/NekoBoxForAndroid/releases")),
		AndroidList.Row(AndroidListBackBtn),
	)

	WindowsList.Inline(
		WindowsList.Row(WindowsList.URL("Hiddify-next (github)", "https://github.com/hiddify/hiddify-next/releases")),
		WindowsList.Row(WindowsListBackBtn),
	)

	MacOSList.Inline(
		MacOSList.Row(MacOSList.URL("Hiddify-next (github)", "https://github.com/hiddify/hiddify-next/releases")),
		MacOSList.Row(MacOSListBackBtn),
	)
}