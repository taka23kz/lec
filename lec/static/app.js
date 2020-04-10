const app = new Vue({
    el: '#app',
    data: {
      questionID: '',
      question: '',
      choices: [],
      choiceNum: 4,
      answerType : '',
      ownerGroupID: 1,
      lessonID: 1,
    },
    methods: {
      regist: () => {
        const payload = {
          'question': app.question, 
          'choices': app.choices,
          'choiceNum': app.choiceNum, 
          'answerType' : app.answerType, 
          'ownerGroupID' : app.ownerGroupID, 
          'lessonID' : app.lessonID
        }
        axios.post('/api/question/regist', payload)
          .then(() => {
            app.question = ''
            app.choice = []
            app.answerType = ''
          })
          .catch((err) => {
            alert(err.response.data.error)
          })
      },
      update: () => {
        const payload = {'question_id' : app.question_id, 'choices' : app.choices, 'answer_type' : app.answer_type}
        axios.post('/api/question/update', payload)
          .then(() => {
            app.question = ''
            app.choice = []
            app.answer_type = ''
          })
          .catch((err) => {
            alert(err.response.data.error)
          })
      },
      select: () => {
        const payload = {'ownerGroupID' : app.ownerGroupID}
        axios.post('api/question/select', payload)
        .then((response) => {
          alert(JSON.stringify(response))
        })
        .catch((err) => {
          alert(err.response.data.error)
        })
      },
      quiz: () => {
        const payload = {'ownerGroupID' : app.ownerGroupID}
        axios.post('api/question/quiz', payload)
        .then((response) => {
          alert(JSON.stringify(response))
        })
        .catch((err) => {
          alert(err.response.data.error)
        })
      }
    }
  })