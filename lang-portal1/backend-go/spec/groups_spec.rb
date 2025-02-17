require 'rspec'
require 'httparty'
require 'json'

def api_url
  'http://localhost:8080/api'
end

RSpec.describe 'Groups API' do
  let(:base_url) { "#{api_url}/groups" }
  let(:valid_group) do
    {
      name: 'JLPT N5 Vocabulary',
      description: 'Essential vocabulary for JLPT N5 level'
    }
  end

  describe 'GET /groups' do
    before do
      @response = HTTParty.get(base_url)
    end

    it 'returns 200 status code' do
      expect(@response.code).to eq(200)
    end

    it 'returns a paginated list of groups' do
      expect(@response.parsed_response).to include('items', 'current_page', 'total_pages', 'total_items', 'items_per_page')
    end

    it 'returns items as an array' do
      expect(@response.parsed_response['items']).to be_an(Array)
    end
  end

  describe 'POST /groups' do
    context 'with valid parameters' do
      before do
        @response = HTTParty.post(
          base_url,
          body: valid_group.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 201 status code' do
        expect(@response.code).to eq(201)
      end

      it 'returns the created group' do
        expect(@response.parsed_response).to include(
          'id' => be_a(Integer),
          'name' => valid_group[:name]
        )
      end
    end

    context 'with invalid parameters' do
      before do
        @response = HTTParty.post(
          base_url,
          body: { name: '' }.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 400 status code' do
        expect(@response.code).to eq(400)
      end

      it 'returns a 404 message' do
        expect(@response.parsed_response).to eq('404 page not found')
      end
    end
  end

  describe 'GET /groups/:id' do
    context 'when group exists' do
      before do
        # Create a group first
        post_response = HTTParty.post(
          base_url,
          body: valid_group.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
        @group_id = post_response.parsed_response['id']
        @response = HTTParty.get("#{base_url}/#{@group_id}")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'returns the group details' do
        expect(@response.parsed_response).to include(
          'id' => @group_id,
          'name' => valid_group[:name],
          'study_session_count' => be_a(Integer),
          'success_rate' => be_a(Integer),
          'word_count' => be_a(Integer)
        )
      end
    end

    context 'when group does not exist' do
      before do
        @response = HTTParty.get("#{base_url}/999999")
      end

      it 'returns 404 status code' do
        expect(@response.code).to eq(404)
      end

      it 'returns a 404 message' do
        expect(@response.parsed_response).to eq('404 page not found')
      end
    end
  end

  describe 'PUT /groups/:id' do
    before do
      # Create a group first
      post_response = HTTParty.post(
        base_url,
        body: valid_group.to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @group_id = post_response.parsed_response['id']
    end

    context 'with valid parameters' do
      before do
        @updated_group = valid_group.merge(name: 'Updated JLPT N5')
        @response = HTTParty.put(
          "#{base_url}/#{@group_id}",
          body: @updated_group.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'returns the updated group' do
        expect(@response.parsed_response).to include(
          'id' => @group_id,
          'name' => 'Updated JLPT N5'
        )
      end
    end

    context 'with invalid parameters' do
      before do
        @response = HTTParty.put(
          "#{base_url}/#{@group_id}",
          body: { name: '' }.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 400 status code' do
        expect(@response.code).to eq(400)
      end

      it 'returns a 404 message' do
        expect(@response.parsed_response).to eq('404 page not found')
      end
    end
  end

  describe 'DELETE /groups/:id' do
    context 'when group exists' do
      before do
        # Create a group first
        post_response = HTTParty.post(
          base_url,
          body: valid_group.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
        @group_id = post_response.parsed_response['id']
        @response = HTTParty.delete("#{base_url}/#{@group_id}")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'removes the group' do
        get_response = HTTParty.get("#{base_url}/#{@group_id}")
        expect(get_response.code).to eq(404)
      end
    end

    context 'when group does not exist' do
      before do
        @response = HTTParty.delete("#{base_url}/999999")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end
    end
  end

  describe 'POST /groups/:id/words/:wordId' do
    before do
      # Create a group first
      post_response = HTTParty.post(
        base_url,
        body: valid_group.to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @group_id = post_response.parsed_response['id']

      # Create a word
      word_response = HTTParty.post(
        "#{api_url}/words",
        body: {
          japanese: '猫',
          romaji: 'neko',
          english: 'cat',
          parts: ['noun']
        }.to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @word_id = word_response.parsed_response['id']
    end

    context 'when both group and word exist' do
      before do
        @response = HTTParty.post("#{base_url}/#{@group_id}/words/#{@word_id}")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'adds the word to the group' do
        get_response = HTTParty.get("#{base_url}/#{@group_id}")
        expect(get_response.parsed_response).to include(
          'id' => @group_id,
          'name' => valid_group[:name]
        )
      end
    end

    context 'when group does not exist' do
      before do
        @response = HTTParty.post("#{base_url}/999999/words/#{@word_id}")
      end

      it 'returns 404 status code' do
        expect(@response.code).to eq(404)
      end

      it 'returns a 404 message' do
        expect(@response.parsed_response).to eq('404 page not found')
      end
    end

    context 'when word does not exist' do
      before do
        @response = HTTParty.post("#{base_url}/#{@group_id}/words/999999")
      end

      it 'returns 404 status code' do
        expect(@response.code).to eq(404)
      end

      it 'returns a 404 message' do
        expect(@response.parsed_response).to eq('404 page not found')
      end
    end
  end

  describe 'DELETE /groups/:id/words/:wordId' do
    before do
      # Create a group
      post_response = HTTParty.post(
        base_url,
        body: valid_group.to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @group_id = post_response.parsed_response['id']

      # Create a word
      word_response = HTTParty.post(
        "#{api_url}/words",
        body: {
          japanese: '猫',
          romaji: 'neko',
          english: 'cat',
          parts: ['noun']
        }.to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @word_id = word_response.parsed_response['id']

      # Add word to group
      HTTParty.post("#{base_url}/#{@group_id}/words/#{@word_id}")
    end

    context 'when both group and word exist' do
      before do
        @response = HTTParty.delete("#{base_url}/#{@group_id}/words/#{@word_id}")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'removes the word from the group' do
        get_response = HTTParty.get("#{base_url}/#{@group_id}")
        word_count = get_response.parsed_response['word_count']
        expect(word_count).to eq(0)
      end
    end

    context 'when group does not exist' do
      before do
        @response = HTTParty.delete("#{base_url}/999999/words/#{@word_id}")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end
    end

    context 'when word does not exist' do
      before do
        @response = HTTParty.delete("#{base_url}/#{@group_id}/words/999999")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end
    end
  end
end
