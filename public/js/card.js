Vue.component("task-card", {
    template: `
  <div class="bg-white shadow rounded px-3 pt-3 pb-5 border border-white">
    <div class="flex justify-between">
      <p class="text-gray-700 font-semibold font-sans tracking-wide text-sm">{{task.title}}</p>
        <!--
      <img
        class="w-6 h-6 rounded-full ml-3"
        src="https://pickaface.net/gallery/avatar/unr_sample_161118_2054_ynlrg.png"
        alt="Avatar"
      >
      -->
    </div>
    <div class="flex mt-4 justify-between items-center">
      <span class="text-sm text-gray-600">{{task.date}}</span>
      <div v-if="task.type" 
        class="px-3 h-6 rounded-full text-xs font-semibold flex items-center"
        :class=""
        >
        <span class="w-2 h-2 rounded-full mr-1" :class=""></span>
        {{task.type}}
        </div>
    </div>
  </div>
  `,
    props: {
        task: {
            type: Object,
            default: () => ({})
        }
    },
    computed: {
        color() {
            const mappings = {
                Design: "purple",
                "Feature Request": "teal",
                Backend: "blue",
                QA: "green",
                default: "teal"
            };
            return mappings[this.task.type] || mappings.default;
            //`bg-${color}-100 text-${color}-700`
            //`bg-${color}-400`
        }
    }
});
