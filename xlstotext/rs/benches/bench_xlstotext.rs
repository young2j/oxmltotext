use criterion::{criterion_group, criterion_main, Criterion};
use xlstotext::bench_run;

fn benchmark(c: &mut Criterion) {
    let path = "../../filesamples/file-sample_100kb.xls";
    c.bench_function("xlstotext", |b| b.iter(|| bench_run(path)));
}

criterion_group!(benches, benchmark);
criterion_main!(benches);
