const quiz = new Vue({
    el: '#quiz',
    data: {
      ownerGroupID: 1,
      question : {
        questionId : '',
        question : '',
        answerType : '',
        choiceNum : '',
        ownerGroupId : '',
        lessonId : '',
      },
      choiceId : '',
      choiceIds : [],
      choices : []
    },
    methods: {
      quiz: () => {
        const payload = {'ownerGroupID' : quiz.ownerGroupID}
        axios.post('api/question/quiz', payload)
        .then((response) => {
          //alert(JSON.stringify(response));
          //alert(response.data.Question.question);
          quiz.question.questionId = response.data.Question.questionId;
          quiz.question.question = response.data.Question.question;
          quiz.question.answerType = response.data.Question.answerType;

          response.data.Choices.forEach(function(v,i) {
            quiz.choices[i] = new Object();
            quiz.choices[i].choiceId = v.choiceId;
            quiz.choices[i].questionId = v.questionId;
            quiz.choices[i].choiceLabel = v.choiceLabel;
            quiz.choices[i].correct = v.correct;
          });
        })
        .catch((err) => {
          alert(err.response.data.error)
        })
      },
      window:onload = function() {
        quiz.quiz();
      },
      answer: () => {
        const payload = {
          'questionID' : quiz.question.questionId,
          'choice' : quiz.choiceId
        }
        axios.post('api/question/answer', payload)
        .then((response) => {
          if ( response.data.correct ) {
            alert('正解！');
          } else {
            alert('不正解!');
          }
          //alert(response.data.answer);
        }).catch((err) => {
          alert(err.response.data.error)
        })
      },
    }
  })