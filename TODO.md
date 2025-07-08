# Bahá'í Transliterator - TODO and Status

## Current Status

### ✅ Completed
- **External Dictionary System**: Successfully moved from hardcoded mappings to external JSON files
- **Modular Architecture**: Clean separation between data and code
- **First Test Case Passing**: Arabic Test Case 1 (Short Prayer) now passes completely
- **Comprehensive Dictionaries**: Created detailed Arabic and Persian dictionaries with:
  - Common words with categories and metadata
  - Divine names and attributes
  - Multi-word phrases
  - Vowel patterns and morphological rules
  - Ezafe rules for Persian
  - Article combination rules for Arabic
- **Fallback System**: Robust fallback when dictionary files aren't available
- **Post-processing Pipeline**: Systematic cleanup and formatting rules

### 🔄 In Progress
- **Expanding Test Coverage**: Working through remaining 3 Arabic test cases
- **Dictionary Completion**: Adding missing words identified from test failures

### ❌ Remaining Work

## Phase 1: Complete Test Suite (Immediate)

### Arabic Tests (3 remaining)
1. **Short Prayer (Test Case 1)** - Minor issues:
   - Missing comma after "Iláhí" 
   - Complex divine name combination: need "anta'l-Mu'ṭí'l-'Alímu'l-Ḥakím"

2. **Prayer of Gratitude (Test Case 3)** - Issues:
   - Complex divine name combinations
   - Article contractions need work
   - Vowel insertion heuristics need refinement
   - Capitalization issues

3. **Lawḥ-i-Ahmad Opening (Test Case 4)** - Issues:
   - Parentheses formatting `*(Lawḥ-i-Aḥmad)*`
   - Complex compound words
   - Proper noun handling
   - Heuristic transliteration producing garbage

4. **Prayer for Purification (Test Case 5)** - Issues:
   - Line break handling
   - Long compound phrases
   - Divine attribute combinations
   - Complex morphological patterns

### Persian Tests (Not yet attempted)
- All 5 Persian test cases need work
- Ezafe rule implementation
- Persian-specific morphology
- Mixed Arabic-Persian vocabulary

## Phase 2: Dictionary Enhancement

### Missing Critical Words
From failing tests, need to add:
- All verb forms with proper voweling
- Compound divine names
- Complex preposition + article combinations
- Morphological variations (plurals, possessives, verb conjugations)

### Improve Heuristics
- **Vowel Insertion**: Better rules for missing vowels
- **Consonant Clusters**: Handle complex combinations
- **Syllable Patterns**: Implement proper syllabification
- **Stress Patterns**: Handle word stress correctly

## Phase 3: Database Testing

### Sample Testing (Next Priority)
1. **Select 10 Arabic prayers** from actual database
2. **Run transliterator** and evaluate quality
3. **Add failing cases** as new test cases
4. **Iterate improvements** until acceptable quality
5. **Repeat with 10 Persian prayers**
6. **Second round** with different 10 prayers each language

### Quality Metrics
- Define what "acceptable quality" means
- Create evaluation criteria
- Document improvement methodology

## Phase 4: Production Deployment

### Full Database Processing
1. **Backup existing transliterations**
2. **Process all Arabic prayers** → regenerate `ar-translit` field
3. **Process all Persian prayers** → regenerate `fa-translit` field
4. **Validate results** with spot checks
5. **Document changes** and improvements

## Technical Debt and Improvements

### Code Quality
- [ ] Add proper error handling throughout
- [ ] Improve logging and debugging
- [ ] Add more comprehensive unit tests
- [ ] Document all functions and data structures
- [ ] Add CLI argument validation

### Performance
- [ ] Profile dictionary loading performance
- [ ] Optimize regex compilation (do once at startup)
- [ ] Consider caching for repeated words
- [ ] Benchmark with large texts

### Maintenance
- [ ] Create dictionary update workflow
- [ ] Add word frequency analysis tools
- [ ] Create validation tools for dictionaries
- [ ] Document dictionary structure and conventions

## Architecture Decisions Made

### ✅ Good Decisions
- **External dictionaries**: Much more maintainable than hardcoded
- **Fallback system**: Tool works even without dictionary files
- **Structured data**: Categories, metadata, and linguistic information
- **Test-driven approach**: Real examples driving implementation

### 🔄 Decisions to Revisit
- **File paths**: Dictionary loading could be more robust
- **Phrase tokenization**: Current system works but could be more elegant
- **Post-processing order**: May need optimization for complex cases

## Resources Needed

### Linguistic Expertise
- Native Arabic/Persian speakers for quality validation
- Bahá'í transliteration experts for convention verification
- Test case expansion with edge cases

### Technical Infrastructure
- Database backup and restore procedures
- Automated testing pipeline
- Performance monitoring tools

## Success Criteria

### Phase 1 Complete When:
- All 10 test cases pass (5 Arabic + 5 Persian)
- Comprehensive dictionary coverage for test vocabulary
- Robust handling of edge cases in tests

### Phase 2 Complete When:
- 90%+ quality on 20 real prayers (10 Arabic + 10 Persian)
- Documented quality metrics and evaluation process
- Stable, well-tested codebase

### Ready for Production When:
- Quality metrics consistently met
- Full database backup procedures in place
- Rollback plan documented and tested
- Stakeholder approval obtained

## Current Status Update

### ✅ Major Fixes Completed
- **Line break preservation**: Fixed regex in post-processing that was removing newlines
- **Phrase token handling**: Fixed punctuation extraction that was breaking phrase replacement
- **Dictionary word variants**: Added missing diacritic variants (وغنآئك, القيّوم)
- **Test Case 2 complete**: All major issues resolved

### 🔄 Current Blockers

1. **Divine name combinations**: Need better compound phrase handling for complex divine names
2. **Missing vocabulary**: Several compound words still missing from dictionary
3. **Heuristic transliteration**: Fallback system producing poor results for unknown words
4. **Persian transliteration**: Fundamentally broken, needs complete revision

## Next Session Priority

1. Fix Test Case 1 (minor comma and divine name issues)
2. Systematically work through Tests 3-5 for Arabic
3. Begin comprehensive Persian test implementation
4. Start planning database sampling strategy

---

*Last updated: 2024-12-19*
*Status: Phase 1 in progress - 2/10 tests passing*