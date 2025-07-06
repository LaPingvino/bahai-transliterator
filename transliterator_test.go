package transliterator

import (
	"testing"
)

func TestArabicTransliteration(t *testing.T) {
	trans := New()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Short Prayer - Test Case 1",
			input:    "يا إِلهِي اسْمُكَ شِفائِي وَذِكْرُكَ دَوائِي وَقُرْبُكَ رَجَائِيْ وَحُبُّكَ مُؤْنِسِيْ وَرَحْمَتُكَ طَبِيبِيْ وَمُعِيْنِيْ فِي الدُّنْيا وَالآخِرَةِ وَإِنَّكَ أَنْتَ المُعْطِ العَلِيمُ الحَكِيمُ.",
			expected: "Yá Iláhí, ismuka shifá'í wa-dhikruka dawá'í wa-qurbuka rajá'í wa-ḥubbuka mu'nisí wa-raḥmatuka ṭabíbí wa-mu'íní fí'd-dunyá wa'l-ákhirati wa-innaka anta'l-Mu'ṭí'l-'Alímu'l-Ḥakím.",
		},
		{
			name:     "Short Obligatory Prayer - Test Case 2",
			input:    "# إِلهِي إِلهِي\nأشهد يا إلهي بأنّك خلقتني لعرفانك وعبادتك. أشهد في هذا الحين بعجزي وقوّتك وضعفي واقتدارك وفقري وغنآئك. لا إله إلاّ أنت المهيمن القيّوم.",
			expected: "# Iláhí Iláhí\nAshhadu yá Iláhí bi-annaka khalaqtaní li-'irfánika wa-'ibádatika. Ashhadu fí hádhá'l-ḥíni bi-'ajzí wa-quwwatika wa-ḍa'fí wa-iqtidárika wa-faqrí wa-ghaná'ika. Lá iláha illá anta'l-Muhayminu'l-Qayyúm.",
		},
		{
			name:     "Prayer of Gratitude - Test Case 3",
			input:    "# هُوَ اللهُ تَعَالى شأنُهُ العَظَمَةُ والاقْتِدارُ\nإِلهِي إِلهِي، أَشْكُرُكَ فِي كُلِّ حالٍ وَأَحْمَدُكَ فِي جَمِيعِ الأَحْوالِ. فِي النِّعْمَةِ أَلْحَمْدُ لَكَ يا إِلهَ العَالَمِينَ. وَفِي فَقْدِها الشُّكْرُ لَكَ يا مَقْصُودَ العَارِفينَ.",
			expected: "# Huwa'lláhu Ta'álá Sha'nuhu'l-'Aẓamatu wa'l-Iqtidár\nIláhí Iláhí, ashkuruka fí kulli ḥálin wa-aḥmaduka fí jamí'i'l-aḥwál. Fí'n-ni'mati al-ḥamdu laka yá Iláha'l-'álamín. Wa-fí faqdihá'sh-shukru laka yá Maqṣúda'l-'árifín.",
		},
		{
			name:     "Lawh-i-Ahmad Opening - Test Case 4",
			input:    "* (لوح احمد)\n# هُوَ السُّلْطَانُ العَليْمُ الحَكِيمُ\nهَذِهِ وَرْقَةُ الفِردَوْسِ تُغَنِّي عَلَى أَفْنَانِ سِدْرَةِ البَقاءِ بِأَلْحانِ قُدْسٍ مَلِيحٍ وتُبَشِّرُ المُخْلِصِينَ إِلَى جِوارِ اللهِ وَالمُوَحِّدِينَ إِلى سَاحَةِ قُرْبٍ كَرَيمٍ",
			expected: "*(Lawḥ-i-Aḥmad)*\n# Huwa's-Sulṭánu'l-'Alímu'l-Ḥakím\nHádhihi waraqatu'l-Firdawsi tughanní 'alá afnáni sidrati'l-baqá'i bi-alḥáni qudsin malíḥin wa-tubashshshiru'l-mukhlliṣína ilá jiwári'lláhi wa'l-muwaḥḥidína ilá sáḥati qurbin karím",
		},
		{
			name:     "Prayer for Purification - Test Case 5",
			input:    "# بِسْمِهِ المُهَيْمِنِ عَلَى الأَسْماءِ\nإِلهِي إِلهِي أَسْأَلُكَ بِبَحْرِ شِفَائِكَ وإِشْراقَاتِ أنْوَارِ نَيِّرِ فَضْلِكَ وَبِالاسْمِ الَّذِي سَخَّرْتَ بِهِ عِبَادَكَ وبِنُفُوذِ كَلِمَتكَ العُلْيَا واقْتِدارِ قَلَمِكَ الأَعْلَى",
			expected: "# Bismihi'l-Muhaymini 'alá'l-Asmá'\nIláhí Iláhí! As'aluka bi-baḥri shifá'ika wa-ishráqáti anwári nayyiri faḍlika wa-bi'l-ismi'lladhí sakhkharta bihi 'ibádaka wa-bi-nufúdhi kalimatika'l-'ulyá wa-iqtidári qalamika'l-a'lá",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trans.Transliterate(tt.input, Arabic)
			if result != tt.expected {
				t.Errorf("Arabic transliteration failed for %s\nInput: %s\nExpected: %s\nGot: %s", 
					tt.name, tt.input, tt.expected, result)
			}
		})
	}
}

func TestPersianTransliteration(t *testing.T) {
	trans := New()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Prayer for Guidance - Test Case 1",
			input:    "اِلهَا مَعبُودا مَلِكا مَلِك اَلمُلُوكا\nاز تو مي‌طلبم تأييد فرمائی و توفيق عطا كنی تا به آنچه سزاوارِ ايّام تو است عمل نمايم و قابلِ جود و كرمِ تو است مشغول گردم. ای كريم غافلان را به‌بحرِ آگاهی راه نما و كنيزانت را به‌انوارِ اسمت مُنَوّر فرما و به‌اعمالِ طيّبه طاهره و اخلاق مرضيّه مؤيّد دار. لَكَ الحَمدُ واَلثَّناءُ وَلَكَ الفَضلُ وَالعَطاءُ. اين نَملۀ فانيه را به‌سُرادقِ عِرفانت راه نمودی و در ظلّ خِباء مَجدَت مأوی دادی. توئی بخشنده و توانا و دانا و بينا.",
			expected: "Iláhá ma'búdan malikan malika'l-mulúk\nAz tú mí-ṭalabam ta'yíd farmá'í va tawfíq 'aṭá kuní tá bih ánchih sizávár-i ayyám-i tú ast 'amal namáyam va qábil-i júd va karam-i tú ast mashghúl gardam. Ay karím! Gháfilán rá bih-baḥr-i ágáhí ráh namá va kanízánat rá bih-anwár-i ismat munawwar farmá va bih-a'mál-i ṭayyibih ṭáhirih va akhlá-i marḍíyyih mu'ayyad dár. Laka'l-ḥamdu wa'th-thaná'u wa-laka'l-faḍlu wa'l-'aṭá'. Ín namliy-i fáníh rá bih-surádiq-i 'irfánat ráh namúdí va dar ẓill-i khíbá'-i majdat ma'vá dádí. Tú'í bakhshindih va tavána va dána va bíná.",
		},
		{
			name:     "Prayer of Supplication - Test Case 2",
			input:    "اِلهَا مَعبُودا مَقصُودا\nفقيری از فُقَراء قصدِ بحرِ عطا نموده و جاهلی از جُهلاء به تجلّياتِ آفتابِ علمت توجّه كرده. سؤال مي‌كنم تو را به دِمائی كه در راهِ تو در ايران ريخته شد و به نفوسی كه سَطوتِ ظالمين و ظلمِ مشركين ايشان‌را از توجّه به تو منع ننمود و از تقَرّب باز نداشت، اينكه كنيزِ خود را از نُعاقِ ناعقين و شُبهاتِ مُريبين حفظ فرمائی و در ظِلّ قِبابِ اسم كريمت مأوی دهی. توئی قادر بر كُلّ و مُهيمِن بر كلّ. اَشهَدُ وَ تَشهَدُ الأشياءُ كُلّها بِاَنَّكَ اَنتَ المُقتَدِرُالقَديرُ.",
			expected: "Iláhá ma'búdan maqṣúd\nFaqírí az fuqará qasad-i baḥr-i 'aṭá namúdih va jáhilí az juhalá bih tajallíyát-i áftáb-i 'ilmat tavahjuh kardih. Su'ál mí-kunam tú rá bih dimá'í kih dar ráh-i tú dar Írán rískhtih shud va bih nufúsí kih saṭvat-i ẓálimín va ẓulm-i mushrikín íshán rá az tavahjuh bih tú man' nanmúd va az taqarrub báz nadásht, ínkih kanīz-i khud rá az nu'áq-i ná'iqín va shubahát-i muríbín ḥifẓ farmá'í va dar ẓill-i qibáb-i ism-i karīmat ma'vá dihí. Tú'í qádir bar kull va muhaymin bar kull. Ashhadu va tash-hadu'l-ashyá'u kulluhá bi-annaka anta'l-Muqtadiru'l-Qadír.",
		},
		{
			name:     "Prayer with Divine Invocation - Test Case 3",
			input:    "بِسْمه المُهيمن القَيُّوم\n#\"\"ای كنيزِ من، به‌اين بيان كه از مَشرِقِ فَمِ رحمن اشراق نموده ناطق باش\"\"\nای پروردگارِ من و يكتا خداوندِ بي‌مانندِ من، شهادت مي‌دهم به يكتائی تو و به‌اينكه از برای تو وزير و معينی نبوده و نيست. لَم يَزَل يكتا بوده‌ای و لا يَزال خواهی بود. ای خدایِ من و محبوبِ جانِ من، امروز روزي‌است كه فُراتِ رحمت جاری و آفتابِ كَرَم مُشرِق و سماءِ عنايت مُرتَفَع است.",
			expected: "Bismihi'l-Muhaymini'l-Qayyúm\n# \"Ay kanīz-i man, bih-ín bayán kih az mashriq-i fam-i Raḥmán ishráq namúdih náṭiq básh\"\nAy Parvardigár-i man va yaktá Khudávand-i bí-mánand-i man, shahádaat mí-diham bih yaktá'í-i tú va bih-ínkih az barí-yi tú vazír va mu'íní nabúdih va níst. Lam yazal yaktá búdih-í va lá yazál kháhí búd. Ay Khudí-yi man va maḥbúb-i ján-i man, imrúz rúzí-st kih Furát-i raḥmat járí va áftáb-i karam mushriq va samá'-i 'ináyat murtafi' ast.",
		},
		{
			name:     "Prayer of Witnessing - Test Case 4",
			input:    "*(بگو ای الهِ من و محبوبِ من و سيّدِ من و سَنَدِ من و مقصودِ من)\n\nشهادت مي‌دهد جان و روان و لسان به‌اينكه واحد بوده‌ای و شبه و مثل برای تو نبوده و نخواهد بود. جميع السُن از وصفت عاجز و جميعِ قلوب از عِرفانت قاصر بوده و خواهد بود. ای پروردگار من، عجز و فقر و فنای كنيزِ خود را مشاهده مي‌نمائی. اين سائلی است ارادۀ بابِ تو نموده و فقيری است قصدِ دريای غَنای تو كرده.",
			expected: "*(Bigú ay Iláh-i man va maḥbúb-i man va sayyid-i man va sanad-i man va maqṣúd-i man)*\n\nShahádaat mí-dihad ján va ruvána va lisán bih-ínkih váḥid búdih-í va shabah va mithl barí-yi tú nabúdih va nakhváhad búd. Jamí'u'l-alsun az waṣfat 'ájiz va jamí'-i qulúb az 'irfánat qáṣir búdih va kháhad búd. Ay Parvardigár-i man, 'ajz va faqr va faná-yi kanīz-i khud rá musháhidih mí-namá'í. Ín sá'ilí-st irádiy-i báb-i tú namúdih va faqírí-st qaṣad-i daryí-yi ghaná-yi tú kardih.",
		},
		{
			name:     "Prayer of Unity and Submission - Test Case 5",
			input:    "كَريما رَحيما\nگواهی مي‌دهم به وحدانيّت و فردانيّت تـو و از تـو مي‌طلبم آنچه را كه به‌دوامِ مُلك و مَلَكوت باقی و پاينده است. توئی مالكِ مَلَكوت و سلطانِ غيب و شهود. ای پروردگار، مسكينی به‌بحرِ غنايت توجّه نموده و سائلی به‌ذيلِ كرمت اقبال كرده، او را محروم منما. توئی آن فَضّالی كه ذرّاتِ كائنات بر فَضلت گواهی داده، توئی آن بخشنده‌ای كه جميع مُمكِنات بر بخششت اعتراف نموده.",
			expected: "Karīman Raḥīm\nGuvāhī mī-diham bih vaḥdāniyyat va fardāniyyat-i tū va az tū mī-ṭalabam ānchih rā kih bih-davām-i mulk va malakūt bāqī va pāyindih ast. Tū'ī mālik-i malakūt va sulṭān-i ghayb va shuhūd. Ay Parvardigār, miskīnī bih-baḥr-i ghanāyat tavahjuh namūdih va sā'ilī bih-dhayl-i karamat iqbāl kardih, ū rā maḥrūm manmā. Tū'ī ān faḍḍālī kih dharrāt-i kā'ināt bar faḍlat guvāhī dādih, tū'ī ān bakhshindih-ī kih jamī'-i mumkināt bar bakhshishat i'tirāf namūdih.",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trans.Transliterate(tt.input, Persian)
			if result != tt.expected {
				t.Errorf("Persian transliteration failed for %s\nInput: %s\nExpected: %s\nGot: %s", 
					tt.name, tt.input, tt.expected, result)
			}
		})
	}
}

func TestLanguageDetection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Language
	}{
		{
			name:     "Arabic with emphatic letters",
			input:    "يا إِلهِي اسْمُكَ شِفائِي",
			expected: Arabic,
		},
		{
			name:     "Persian with specific letters",
			input:    "پروردگار چه کنم",
			expected: Persian,
		},
		{
			name:     "Persian with common words",
			input:    "از تو می‌طلبم",
			expected: Persian,
		},
		{
			name:     "Arabic default for mixed",
			input:    "الله أكبر",
			expected: Arabic,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AutoDetectLanguage(tt.input)
			if result != tt.expected {
				t.Errorf("Language detection failed for %s\nInput: %s\nExpected: %v\nGot: %v", 
					tt.name, tt.input, tt.expected, result)
			}
		})
	}
}

func TestCommonPatterns(t *testing.T) {
	trans := New()
	
	tests := []struct {
		name     string
		input    string
		expected string
		lang     Language
	}{
		{
			name:     "Allah name",
			input:    "الله",
			expected: "Alláh",
			lang:     Arabic,
		},
		{
			name:     "Article connection wa-al",
			input:    "والله",
			expected: "wa'lláh",
			lang:     Arabic,
		},
		{
			name:     "No god but formula",
			input:    "لا إله إلا الله",
			expected: "lá iláha illá'lláh",
			lang:     Arabic,
		},
		{
			name:     "Persian God",
			input:    "خدا",
			expected: "Khudá",
			lang:     Persian,
		},
		{
			name:     "Persian Lord",
			input:    "پروردگار",
			expected: "Parvardigár",
			lang:     Persian,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trans.Transliterate(tt.input, tt.lang)
			if result != tt.expected {
				t.Errorf("Pattern test failed for %s\nInput: %s\nExpected: %s\nGot: %s", 
					tt.name, tt.input, tt.expected, result)
			}
		})
	}
}

func BenchmarkArabicTransliteration(b *testing.B) {
	trans := New()
	text := "يا إِلهِي اسْمُكَ شِفائِي وَذِكْرُكَ دَوائِي وَقُرْبُكَ رَجَائِيْ وَحُبُّكَ مُؤْنِسِيْ وَرَحْمَتُكَ طَبِيبِيْ وَمُعِيْنِيْ فِي الدُّنْيا وَالآخِرَةِ وَإِنَّكَ أَنْتَ المُعْطِ العَلِيمُ الحَكِيمُ."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trans.Transliterate(text, Arabic)
	}
}

func BenchmarkPersianTransliteration(b *testing.B) {
	trans := New()
	text := "اِلهَا مَعبُودا مَلِكا مَلِك اَلمُلُوكا از تو مي‌طلبم تأييد فرمائی و توفيق عطا كنی تا به آنچه سزاوارِ ايّام تو است عمل نمايم"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trans.Transliterate(text, Persian)
	}
}