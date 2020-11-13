Vue.component("task-card", {
    template: `
  <div class="bg-white shadow rounded px-3 pt-3 pb-5 border border-white">
    <div class="flex justify-between">
      <p class="text-gray-700 font-semibold font-sans tracking-wide text-sm">{{image.name}}</p>
        <!--
      <img
        class="w-6 h-6 rounded-full ml-3"
        src="https://pickaface.net/gallery/avatar/unr_sample_161118_2054_ynlrg.png"
        alt="Avatar"
      >
      -->
    </div>
    <div class="flex mt-4 justify-between items-center">
      <span  class="text-sm text-gray-600">
      {{image.version | truncate(20, "..") }}
      
      </span>
      <div v-if="image.version" 
        class="px-3 h-6 rounded-full text-xs font-semibold flex items-center"
        :class=""
        >
        <span class="w-2 h-2 rounded-full mr-1" :class=""></span>
        <span v-if="false">{{ (new Date(image.createdOn)).toLocaleString() }}</span>
        </div>
    </div>
  </div>
  `,
    props: {
        image: {
          "name": "unknown",
          "createdOn": "0001-01-01T00:00:00Z",
          "version": "v0.0",
          "status": "New",
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
