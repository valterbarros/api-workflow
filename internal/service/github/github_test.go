package github

import (
	"api-workflow/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLastVersions(t *testing.T) {
	text := utils.GetFileData(utils.GetFromRoot("fixtures/CHANGELOG.md"))

	result := GetLastVersions(text, 3, "VOP")
	expected := `# v1.1413.0
- [VOP-1172] Tratar conversões de documentos excluídos no ConvertWorker [#VOP-1172](https://clicksign.atlassian.net/browse/VOP-1172) > [#14768](https://github.com/clicksign/tavola/pull/14768)
# v1.1412.0
- [VOP-1159] [BACK] Disponibilizar teste A/B em endpoints [#VOP-1159] [BACK](https://clicksign.atlassian.net/browse/VOP-1159] [BACK) > [#14765](https://github.com/clicksign/tavola/pull/14765)
# v1.1411.0
- [VOP-1134] Criar evento de webhook para document ativado ou falha na ativação [#VOP-1134](https://clicksign.atlassian.net/browse/VOP-1134) > [#14754](https://github.com/clicksign/tavola/pull/14754)
- [VOP-1147] (BACK) Implementar filtro por name no endpoint /events conforme JSON:API [#VOP-1147](https://clicksign.atlassian.net/browse/VOP-1147) > [#14750](https://github.com/clicksign/tavola/pull/14750)
- [VOP-1093] Criar endpoint para lidar com a ativação de envelope com muitos documentos [#VOP-1093](https://clicksign.atlassian.net/browse/VOP-1093) > [#14745](https://github.com/clicksign/tavola/pull/14745)
- [VOP-1104] Adiciona botão de não mostrar novamente [#VOP-1104](https://clicksign.atlassian.net/browse/VOP-1104) > [#14742](https://github.com/clicksign/tavola/pull/14742)
# v1.1410.0
`

	if result != expected {
		t.Errorf("error")
	}
}

func TestGetVersion(t *testing.T) {
	text := utils.GetFileData(utils.GetFromRoot("fixtures/CHANGELOG.md"))

	result := GetVersion(text, "VOP-1134")
	expected := `# v1.1411.0
- v1.1410.17 -- [VOP-1134] Criar evento de webhook para document ativado ou falha na ativação [#VOP-1134](https://clicksign.atlassian.net/browse/VOP-1134) > [#14754](https://github.com/clicksign/tavola/pull/14754)`

	if result != expected {
		t.Errorf("error")
	}
}

func TestReviewsResult(t *testing.T) {
	t.Run("when has 16 comments and I approved the pr", func(t *testing.T) {
		ser := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// 16 comments
			// 1 approved
			// 0 dismissed
			data := utils.GetFileData(utils.GetFromRoot("fixtures/reviews.json"))
			w.Write([]byte(data))
		}))

		defer ser.Close()

		result := ReviewsResult(ser.URL)

		if result.color != "yellow" {
			t.Error("not equal color")
		}

		if result.text != "D: CO ME: OK" {
			t.Error("not equal review data")
		}
	})
	t.Run("when has more than 2 approves", func(t *testing.T) {
		ser := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// 14 comments
			// 3 approved
			// 0 dismissed
			// data := utils.GetFileData(utils.GetFromRoot("fixtures/reviews.json"))
			data := utils.GetFileData(utils.GetFromRoot("fixtures/reviews-approved.json"))
			w.Write([]byte(data))
		}))

		defer ser.Close()

		result := ReviewsResult(ser.URL)

		if result.color != "green" {
			t.Error("not equal color")
		}

		if result.text != "D: CO ME: OK" {
			t.Error("not equal review data")
		}
	})

	t.Run("when has no reviews", func(t *testing.T) {
		ser := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// 14 comments
			// 3 approved
			// 0 dismissed
			// data := utils.GetFileData(utils.GetFromRoot("fixtures/reviews.json"))
			w.Write([]byte("[]"))
		}))

		defer ser.Close()

		result := ReviewsResult(ser.URL)

		if result.color != "yellow" {
			t.Error("not equal color")
		}

		if result.text != "No reviews" {
			t.Error("not equal review data")
		}
	})

	t.Run("when has no reviews", func(t *testing.T) {
		ser := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// 14 comments
			// 3 approved
			// 0 dismissed
			// data := utils.GetFileData(utils.GetFromRoot("fixtures/reviews.json"))
			w.Write([]byte(`[{"state": "APPROVED", "user": {"login": "valter-clk"}}]`))
		}))

		defer ser.Close()

		result := ReviewsResult(ser.URL)

		if result.color != "yellow" {
			t.Error("not equal color")
		}

		if result.text != "ME: OK" {
			t.Error("not equal review data result:", result.text)
		}
	})
}
