import { alias } from '@ember/object/computed';
import Controller from '@ember/controller';
import { computed, observer } from '@ember/object';
import { task } from 'ember-concurrency';
import Sortable from 'nomad-ui/mixins/sortable';
import Searchable from 'nomad-ui/mixins/searchable';
import messageFromAdapterError from 'nomad-ui/utils/message-from-adapter-error';

export default Controller.extend(Sortable, Searchable, {
  queryParams: {
    currentPage: 'page',
    searchTerm: 'search',
    sortProperty: 'sort',
    sortDescending: 'desc',
    onlyPreemptions: 'preemptions',
  },

  // Set in the route
  flagAsDraining: false,

  currentPage: 1,
  pageSize: 8,

  sortProperty: 'modifyIndex',
  sortDescending: true,

  searchProps: computed(() => ['shortId', 'name']),

  onlyPreemptions: false,

  visibleAllocations: computed(
    'model.allocations.[]',
    'preemptions.[]',
    'onlyPreemptions',
    function() {
      return this.onlyPreemptions ? this.preemptions : this.model.allocations;
    }
  ),

  listToSort: alias('visibleAllocations'),
  listToSearch: alias('listSorted'),
  sortedAllocations: alias('listSearched'),

  eligibilityError: null,
  showDrainNotification: false,

  preemptions: computed('model.allocations.@each.wasPreempted', function() {
    return this.model.allocations.filterBy('wasPreempted');
  }),

  sortedEvents: computed('model.events.@each.time', function() {
    return this.get('model.events')
      .sortBy('time')
      .reverse();
  }),

  sortedDrivers: computed('model.drivers.@each.name', function() {
    return this.get('model.drivers').sortBy('name');
  }),

  setEligibility: task(function*(value) {
    try {
      yield value ? this.model.setEligible() : this.model.setIneligible();
    } catch (err) {
      const error = messageFromAdapterError(err) || 'Could not set eligibility';
      this.set('eligibilityError', error);
    }
  }).drop(),

  triggerDrainNotification: observer('model.isDraining', function() {
    if (!this.model.isDraining && this.flagAsDraining) {
      this.set('showDrainNotification', true);
    }

    this.set('flagAsDraining', this.model.isDraining);
  }),

  actions: {
    gotoAllocation(allocation) {
      this.transitionToRoute('allocations.allocation', allocation);
    },

    setPreemptionFilter(value) {
      this.set('onlyPreemptions', value);
    },
  },
});
