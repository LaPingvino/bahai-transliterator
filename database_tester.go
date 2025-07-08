package transliterator

import (
	"encoding/json"
	"fmt"
	"os"
)

// DatabaseTestCase represents a test case from the database
type DatabaseTestCase struct {
	SourceID         string `json:"source_id"`
	Language         string `json:"language"`
	Original         string `json:"original"`
	CurrentTranslit  string `json:"current_translit"`
	NewTranslit      string `json:"new_translit"`
	Improved         bool   `json:"improved"`
}

// Test samples from the database
var persianSamples = []DatabaseTestCase{
	{
		SourceID: "1544",
		Language: "fa",
		Original: "هُواللّه\nای يزدانِ مهربان، سرا پا گُنَهيم و خاکِ رَهيم و مُتِضَرّع در هر صبحگهيم، ای بزرگوار، خطا بپوش و عطا ببخش، وفا بفرما، صفا عنايت کن تا نورِ هدايت تابد و پرتوِ موهبت بيفزايد، شمعِ غفران بر افروزد و پردۀ عِصيان بسوزد، صبحِ اميد دَمَد، ظلمتِ نوميد زائل گردد، نسيمِ الطاف بِوَزَد و شميمِ اِحسان مُرور نمايد، مشام‌ها مُعَطّر گردد، روي‌ها مُنَوّر شود. توئی  بخشنده و مهربان و درخشنده و تابان.",
		CurrentTranslit: "Húvállh ay yzdáni mhrbán, sra pa gúnáhym va kháki ráhym va mútídár dár hr sbhghym, ay bzrgvár, khta bpvsh va tá bbkhsh, vfa bfrma, sfa náyt kun tá nvri hdáyt tábd va prtvi mvhbt byfzáyd, shmi ghfrán br afrvzd va prd isyán bsvzd, sbhi amyd dámád, zlmti nvmyd záíl grdd, nsymi altáf bívázád va shmymi aíhsán múrvr nmáyd, mshámha múátr grdd, rvyha múnávr shvd. Tví bkhshndh va mhrbán va drkhshndh va tábán.",
	},
	{
		SourceID: "1543",
		Language: "fa",
		Original: "هُوالابهی\nای خداوندِ مهربان، اين اسيرانِ زنجير مَحبّتت را دستگير شو و مَلجأ و پناه و مُجير، نَفَحاتِ قُدس از گُلشنِ عنايت بفرست و ساحتِ دل‌ها را گُلستانِ موهبت کن و چمنستانِ حقيقت نما، از غيرِ خود بي‌زار کن و به راز و نياز دَمساز فرما، مورِ حَقير را سُليمانِ اِقليمِ جليل کن و ذَرّۀ فقير را اميرِ اوجِ اثير فرما، قطره را موهبتِ بحر بخش و سبزه را طراوت و لطافتِ شجرۀ اَخضَر عطا فـرمـا، کُل  يارانِ تـو انـد و بندگانِ درگاهِ تو، فضل وَ  جود مبذول دار و در اين يومِ مسعود تأييدِ مخصوص مشهود کن. توئی بخشنده و مهربان و درخشنده  و تابان.",
		CurrentTranslit: "Húvúl-Abhá ay khúdávándi mhrbán, ín asyráni znjyr máhbtt rá dstgyr shv va málja va pnáh va mújyr, náfáháti qúds az gúlshni náyt bfrst va sáhti dlha rá gúlstáni mvhbt kun va chmnstáni hqyqt nma, az ghyri khvd byzár kun va bíh ráz va nyáz dámsáz farmá, mvri háqyr rá súlymáni aíqlymi jlyl kun va dhár fqyr rá amyri avji athyr farmá, qtrh rá mvhbti bhr bkhsh va sbzh rá trávt va ltáfti shjr ákhdár tá farmá, kúl yáráni tv and va bndgáni drgáhi tv, fdl va jvd mbdhvl dár va dár ín yvmi msvd táyydi mkhsvs mshhvd kun. Tví bkhshndh va mhrbán va drkhshndh va tábán.",
	},
	{
		SourceID: "1395",
		Language: "fa",
		Original: "هُواللّه\nای پروردگار، به جنودِ مَلأ اَعلی نصرت نما و به جيوشِ ملائکه محبّت و صفا اعانت کن نَغَماتِ قدس بفرست و مَحافِلِ اُنس مُعطّر نما، فيضِ قديم مَبذول دار و فوزِ عظيم شايان نما، نورِ حقيقت جلوه ده ديدۀ اهلِ بَصيرت روشن کن، آهنگِ مَلکوت ابهی به‌گوش رسان و هر دلتنگِ عالمِ اَدنی را خوشوقت کن، ابرِ رحمت بفرست، بارانِ مُوهبت بِبار، چمنِ هدايت بيارا، رياحينِ مَعانی اَنبات کن و سُلطانِ گُل را تاجِ موهبت بر سر نه و بلبلانِ روحانی را به غزلخوانی بخوان و حقايق و معانی تعليم ده. توئی پروردگار توئی کردگار توئی  مجلّی طُور در کشورِ انوار.",
		CurrentTranslit: "Húvállh ay Párvárdígár, bíh jnvdi mála ály nsrt nma va bíh jývshi mláíkh maḥabbat va sfa ánt kun nághámáti qds bfrst va máháfíli aúns mútr nma, fydi qdym mábdhvl dár va fvzi zym sháyán nma, nvri hqyqt jlvh dh dyd ahli básyrt rvshn kun, ahngi málkvt abhy bhgvsh rsán va hr dltngi almi ádny rá khvshvqt kun, abri rhmt bfrst, báráni múvhbt bíbár, chmni hdáyt byára, ryáhyni mány ánbát kun va súltáni gúl rá táji mvhbt br sr nh va blbláni rvhány rá bíh ghzlkhvány bkhván va hqáyq va mány tlym dh. Tví Párvárdígár tví krdgár tví mjly túvr dár kshvri anvár.",
	},
	{
		SourceID: "1496",
		Language: "fa",
		Original: "هواللّه\nخدایا، طفلم در ظلِ عنایتت پرورش ده. نهالِ تازه‌ام به رشحاتِ سحابِ عنایت پرورش فرما. گیاهِ حدیقۀ مَحبتم، درختِ بارور کن. تـوئی مقتدر و توانا و تـوئی مهـربان و دانا و بینا.",
		CurrentTranslit: "Hvállh khdáya, tflm dár zli náytt prvrsh dh. Nháli tázhám bíh rshháti shábi náyt prvrsh farmá. Gyáhi hdyq máhbtm, drkhti bárvr kun. Tví mqtdr va tvána va tví mhrbán va dána va byna.",
	},
	{
		SourceID: "1381",
		Language: "fa",
		Original: "پاكا پادشاها\nهر آگاهی بر يكتائيت گُواهی داده. توئی آن توانائی كه جودت وجود را موجود فرمود و خطای عباد عطايت را باز نداشت. ای كريم از مَطلعِ نورت مُنَوَّر نما و از مشرقِ غَنايت ثروت حقيقی بخش. توئی بخشنده و توانا.",
		CurrentTranslit: "Páka pádsháha hr agáhy br yktáít gúváhy dádh. Tví án tvánáí kh jvdt vjvd rá mvjvd frmvd va khtáy bád táyt rá báz ndásht. Ay krym az mátli nvrt múnávár nma va az mshrqi ghánáyt thrvt hqyqy bkhsh. Tví bkhshndh va tvána.",
	},
}

var arabicSamples = []DatabaseTestCase{
	{
		SourceID: "3266",
		Language: "ar",
		Original: "يا إِلهِي اسْمُكَ شِفائِي وَذِكْرُكَ دَوائِي وَقُرْبُكَ رَجَائِيْ وَحُبُّكَ مُؤْنِسِيْ وَرَحْمَتُكَ طَبِيبِيْ وَمُعِيْنِيْ فِي الدُّنْيا وَالآخِرَةِ وَإِنَّكَ أَنْتَ المُعْطِ العَلِيمُ الحَكِيمُ.",
		CurrentTranslit: "Ya Iláhí asmuka shifaií vadhikruka davaií vaqurbuka rajáií vahubuka múnisí varahmatuka tabíbí vamuíní fí aldunya válakhirahi vaiinaka ánta almuti alalímu alhakímu.",
	},
	{
		SourceID: "3287",
		Language: "ar",
		Original: "# بِسْمِهِ المُهَيْمِنِ عَلَى الأَسْماءِ\nقُلْ إِلهِي إِلهِي، فَرِّجْ هَمِّي بِجُودِكَ وعَطَائِكَ، وأَزِلْ كُرْبَتِي بِسَلْطَنَتِكَ واقْتِدَارِكَ. تَرانِي يا إِلهِي مُقْبِلاً إِليكَ حينَ إِذْ أَحاطَتْ بِيَ الأَحْزَانُ مِنْ كُلِّ الجِّهَاتِ. أَسأَلُكَ يا مَالِكَ الوُجُودِ والمُهَيْمِنَ على الغَيْبِ والشُّهُودِ، باسْمِكَ الَّذي بِهِ سَخَّرْتَ الأَفْئِدَةَ والقُلُوبَ وبِأَمْواجِ بَحْرِ رَحْمَتِكَ وإِشْراقاتِ أَنْوارِ نيِّرِ عَطَائِكَ أَنْ تَجْعَلَنِي مِنَ الَّذينَ ما مَنَعَهُم شَيْءٌ مِنَ الأَشْياءِ عَنْ التَّوجُّهِ إِلَيْكَ يا مَوْلَى الأَسْماءِ وَفاطِرَ السَّماءِ، أَيْ رَبِّ تَرَى ما وَرَدَ عَلَيَّ فِي أَيَّامِكَ، أَسأَلُكَ بِمَشْرِقِ أَسْمَائِكَ ومَطْلِعِ صِفَاتِكَ أنْ تُقَدَّرَ لِي ما يَجْعَلُنِي قَائِمًا على خِدْمَتِكَ وَناطِقًا بِثَنَائِكَ. إِنَّكَ أَنْتَ المُقْتَدِرُ القَدِيرُ وبِالإِجَابَةِ جَدِيرٌ. ثُمَّ أَسْأَلُكَ في آخِرِ عَرْضِي بِأنْوارِ وَجْهِكَ أَنْ تُصْلِحَ أُمُورِي وتَقْضِي دَيْنِي وَحَوائِجِي إِنَّكَ أَنْتَ الّذي شَهِدَ كُلُّ ذِي لِسَانٍ بقُدْرَتِكَ وقُوَّتِكَ، وذِي دِرَايَةٍ بِعَظَمَتِكَ وسُلْطانِكَ. لا إِلهَ إِلاَّ أَنْتَ السَّامِعُ المُجِيبُ.",
		CurrentTranslit: "# Bismihi almuhaymini alá alásmai\nqul Iláhí Iláhí, farij hamí bijuvdika vatáiika, vázil kurbatí bisaltanatika vaqtidárika. Taraní ya Iláhí muqbilán iilyka hyna iidh áhatat biya aláhzánu min kuli aljiháti. ásáluka ya málika alvujuvdi valmuhaymina la alghaybi valshuhuvdi, basmika aladhy bihi sakharta aláfiidaha valquluvba vbiámvaji bahri rahmatika viishraqati ánvari nyiri atáiika án tajalaní mina aladhyna má manáhum shayun mina aláshyai an altavjuhi iilayka ya mavlá alásmai vafatira alsamai, áy Rabbí tará má varada alaya fí áyámika, ásáluka bimashriqi ásmáiika vmatlii sifátika an tuqadara lí má yajaluní qáiimana la khidmatika vanatiqana bithanáiika. Iinaka ánta almuqtadiru alqadíru vbialiijábahi jadírun. Thuma ásáluka fy akhiri ardí bianvari vajhika án tusliha aumuvrí vtaqdí dayní vahavaiijí iinaka ánta aldhy shahida kulu dhí lisánin bqudratika vquvatika, vdhí diráyahin biazamatika vsultanika. La iilha iilá ánta alsámiu almujíbu.",
	},
}

// runDatabaseTests tests the transliterator against database samples
func runDatabaseTests() {
	// Initialize transliterator
	t, err := New()
	if err != nil {
		fmt.Printf("Error initializing transliterator: %v\n", err)
		return
	}
	
	fmt.Println("=== Testing Persian Samples ===")
	for i, sample := range persianSamples {
		fmt.Printf("\n--- Persian Sample %d (Source ID: %s) ---\n", i+1, sample.SourceID)
		fmt.Printf("Original:\n%s\n", sample.Original)
		fmt.Printf("Current Translit:\n%s\n", sample.CurrentTranslit)
		
		// Test our transliterator
		newTranslit := t.Transliterate(sample.Original, Persian)
		fmt.Printf("New Translit:\n%s\n", newTranslit)
		
		// Compare quality
		if len(newTranslit) > 0 && newTranslit != sample.CurrentTranslit {
			fmt.Printf("*** IMPROVED: New transliteration differs from current\n")
			sample.NewTranslit = newTranslit
			sample.Improved = true
		} else {
			fmt.Printf("*** SAME: No significant improvement\n")
			sample.NewTranslit = newTranslit
			sample.Improved = false
		}
		
		// Update the sample
		persianSamples[i] = sample
	}
	
	fmt.Println("\n=== Testing Arabic Samples ===")
	for i, sample := range arabicSamples {
		fmt.Printf("\n--- Arabic Sample %d (Source ID: %s) ---\n", i+1, sample.SourceID)
		fmt.Printf("Original:\n%s\n", sample.Original)
		fmt.Printf("Current Translit:\n%s\n", sample.CurrentTranslit)
		
		// Test our transliterator
		newTranslit := t.Transliterate(sample.Original, Arabic)
		fmt.Printf("New Translit:\n%s\n", newTranslit)
		
		// Compare quality
		if len(newTranslit) > 0 && newTranslit != sample.CurrentTranslit {
			fmt.Printf("*** IMPROVED: New transliteration differs from current\n")
			sample.NewTranslit = newTranslit
			sample.Improved = true
		} else {
			fmt.Printf("*** SAME: No significant improvement\n")
			sample.NewTranslit = newTranslit
			sample.Improved = false
		}
		
		// Update the sample
		arabicSamples[i] = sample
	}
}

// saveResultsToJSON saves the test results to a JSON file
func saveResultsToJSON() {
	results := map[string]interface{}{
		"persian_samples": persianSamples,
		"arabic_samples":  arabicSamples,
	}
	
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	
	err = os.WriteFile("database_test_results.json", data, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		return
	}
	
	fmt.Printf("\nResults saved to database_test_results.json\n")
}

// RunDatabaseTests runs the database tests
func RunDatabaseTests() {
	fmt.Println("Bahai Transliterator Database Test")
	fmt.Println("==================================")
	
	runDatabaseTests()
	saveResultsToJSON()
	
	// Print summary
	persianImproved := 0
	arabicImproved := 0
	
	for _, sample := range persianSamples {
		if sample.Improved {
			persianImproved++
		}
	}
	
	for _, sample := range arabicSamples {
		if sample.Improved {
			arabicImproved++
		}
	}
	
	fmt.Printf("\n=== SUMMARY ===\n")
	fmt.Printf("Persian: %d/%d samples improved\n", persianImproved, len(persianSamples))
	fmt.Printf("Arabic: %d/%d samples improved\n", arabicImproved, len(arabicSamples))
	fmt.Printf("Total: %d/%d samples improved\n", persianImproved+arabicImproved, len(persianSamples)+len(arabicSamples))
}