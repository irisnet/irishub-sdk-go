package integration_test

import "strconv"

type TestRandom struct{}

var testRandom TestRandom

func (s IntegrationTestSuite) TestRandom() {

	cases := []SubTest{
		{
			"TestQueryRandom",
			queryRandom,
		},
		{
			"TestQueryRandomRequestQueue",
			queryRandomRequestQueue,
		},
	}

	for _, t := range cases {
		s.Run(t.testName, func() {
			t.testCase(s)
		})
	}
}

func queryRandom(s IntegrationTestSuite) {
	res, err := s.Random.QueryRandom("4fdf899d73b1b5e564aa9fb8afd9ce397de325c77c0633aed4b93c52275197d1")
	s.NoError(err)
	s.NotEmpty(res.RequestTxHash)
	value, _ := strconv.ParseFloat(res.Value, 10)
	s.Greater(value, float64(0))
}

func queryRandomRequestQueue(s IntegrationTestSuite) {
	queue, err := s.Random.QueryRandomRequestQueue(12145)
	s.NoError(err)
	s.NotEmpty(queue)
}
