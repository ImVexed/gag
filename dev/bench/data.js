window.BENCHMARK_DATA = {
  "lastUpdate": 1631504407955,
  "repoUrl": "https://github.com/ImVexed/gag",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "v@nul.lu",
            "name": "V-X",
            "username": "ImVexed"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1269c6bffa8808f02e7458b667f75ffdbcf49e7a",
          "message": "Update benchmark.yml",
          "timestamp": "2021-09-12T20:39:02-07:00",
          "tree_id": "a1c5e40d20cc4d06378b382be32f60f932a4b39a",
          "url": "https://github.com/ImVexed/gag/commit/1269c6bffa8808f02e7458b667f75ffdbcf49e7a"
        },
        "date": 1631504407487,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExampleParallel",
            "value": 1040,
            "unit": "ns/op\t     760 B/op\t       5 allocs/op",
            "extra": "1000000 times\n2 procs"
          },
          {
            "name": "BenchmarkExample",
            "value": 2035,
            "unit": "ns/op\t     760 B/op\t       5 allocs/op",
            "extra": "605318 times\n2 procs"
          },
          {
            "name": "BenchmarkBigExample",
            "value": 20624789,
            "unit": "ns/op\t  247917 B/op\t    2900 allocs/op",
            "extra": "58 times\n2 procs"
          }
        ]
      }
    ]
  }
}