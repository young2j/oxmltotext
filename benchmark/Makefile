bench_join:
	go test -benchmem -bench ^Benchmark_Join -benchtime=1s 

bench_docx:
	go test -benchmem -bench ^Benchmark_ParseDocx -benchtime=1s 

bench_xlsx:
	go test -benchmem -bench ^Benchmark_ParseXlsx -benchtime=1s 

bench_pptx:
	go test -benchmem -bench ^Benchmark_ParsePptx -benchtime=1s 

bench_doc:
	go test -benchmem -bench ^Benchmark_ParseDoc[^x] -benchtime=1s 

bench_xls:
	go test -benchmem -bench ^Benchmark_ParseXls[^x] -benchtime=1s 

bench_ppt:
	go test -benchmem -bench ^Benchmark_ParsePpt[^x] -benchtime=1s 

bench_pdf:
	go test -benchmem -bench ^Benchmark_ParsePdf -benchtime=1s 

bench_all: bench_docx bench_xlsx bench_pptx bench_pdf bench_doc bench_xls bench_ppt
