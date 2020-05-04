const qlist = new Vue({
    el: '#qlist',
    data: {
      // search condition
      questionId : '',      // 検索条件として設定しているquestionId
      ownerGroupId: 1,      // 検索条件として設定しているonwerGroupId（今のところ固定)
      selectLesson : null,
      // selected question list
      lessons : []          // 検索条件として表示されるLessonの一覧(ownerGroupIdに紐づくLessonが設定される)
    },
    methods: {
      init: () => {
        const payload = {
          'ownerGroupID' : qlist.ownerGroupId
        }
        // 出題用の問題を取得し、qlistオブジェクトに格納する。
        axios.post('api/qlist/init', payload)
        .then((response) => {
          //alert(JSON.stringify(response));
          response.data.lessons.forEach(function(v,i) {
            qlist.lessons[i] = { 
              "lessonId" : v.lessonId,
              "lessonName" : v.lessonName
            }
          });
          // ↑でqlist.lessonsの初期化前にselectLessonを設定する画面の初期表示時に選択肢が正しく設定されない。
          // qlist.lessons→qlist.selectLessonの順序で設定する必要がある。
          // 選択した検索条件(lesson)を設定
          qlist.selectLesson = response.data.searchCondition.LessonID;
        })
        .catch((err) => {
          alert(err)
        })
      },
      // 画面ロード時にinitメソッドを実行
      window:onload = function() {
        qlist.init();
      },
    }
  })